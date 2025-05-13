package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

// GetLiquidationPrice Happens when Equity Value = Maintenance Margin Value
// Equity Value = InitialMarginValue + UPnL
// Maintenance Margin Value = MMR * Abs(Qty) * LiqPrice
// UPnL = Qty * (LiqPrice - EntryPrice)
// LiqPrice = (IMV - Qty * Entry Price) / (MMR * Abs(Qty) - Qty)
// We do not simplify this using IMR because IMR can change through gov proposal
// Also this automatically handles Liquidation price when collateral is added
func (k Keeper) GetLiquidationPrice(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket, subAccount types.SubAccount) (math.LegacyDec, error) {
	if subAccount.IsIsolated() {
		quoteDenomPrice, err := k.GetDenomPrice(ctx, market.QuoteDenom)
		if err != nil {
			return math.LegacyZeroDec(), err
		}
		initialMarginValue := perpetual.MarginAmount.ToLegacyDec().Mul(quoteDenomPrice)
		num := initialMarginValue.Sub(perpetual.Quantity.Mul(perpetual.EntryPrice))
		den := (market.MaintenanceMarginRatio.Mul(perpetual.Quantity.Abs())).Sub(perpetual.Quantity)
		// the denominator can be 0 only if MMR is 1, which would mean the position would close instantly
		if den.IsZero() {
			return math.LegacyZeroDec(), errors.New("cannot calculate liquidation price, division by zero (MMR = 1)")
		}
		return num.Quo(den), nil
	} else {
		panic("implement me")
	}
}

// ForcedLiquidation Possible cases in isolated margin:
//  1. Equity Value > 0 and the market absorbs the whole position - liquidator fee paid by margin
//  2. Equity Value > 0 and the market only partially absorbs position - use ADL to close the rest of the position and partial liquidator fee can be paid by margin
//  3. Equity Value <= 0 and there are enough insurance funds, the market absorbs the whole position - liquidator fee paid by insurance funds
//  4. Equity Value <= 0 and there are enough insurance funds, the market partially absorbs the whole position - liquidator fee paid by insurance funds
//  5. Equity Value <= 0 and there are NOT enough insurance funds.
//
// What about slippage for large orders? Can it create problems while executing market orders
// Auto Deleverage: Find positions with maximum profit and close against that.
// If there are not enough Insurance Funds, and we directly delete the position after transferring funds but then this would break the invariant:
// sum of all long = market total open interest == absolute value of sum of all short
func (k Keeper) ForcedLiquidation(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket, liquidator sdk.AccAddress) error {
	subAccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), market.Id)
	if err != nil {
		return fmt.Errorf("forced_liquidation: failed to get subaccount for owner %s, market %d: %w", perpetual.Owner, market.Id, err)
	}
	liquidationPrice, err := k.GetLiquidationPrice(ctx, perpetual, market, subAccount)
	if err != nil {
		return fmt.Errorf("forced_liquidation: failed to get liquidation price for perp %d: %w", perpetual.Id, err)
	}
	currentTwapPrice := k.GetCurrentTwapPrice(ctx, market.Id) // Assuming this is the Mark Price
	if currentTwapPrice.IsZero() {
		return errors.New("forced_liquidation: cannot liquidate, current twap mark price is zero")
	}
	if perpetual.IsLong() && liquidationPrice.LT(currentTwapPrice) { // Price hasn't dropped enough for Long
		return fmt.Errorf("forced_liquidation: long position not yet liquidatable. LiqPx: %s, MarkPx: %s", liquidationPrice.String(), currentTwapPrice.String())
	}
	if perpetual.IsShort() && liquidationPrice.GT(currentTwapPrice) { // Price hasn't risen enough for Short
		return fmt.Errorf("forced_liquidation: short position not yet liquidatable. LiqPx: %s, MarkPx: %s", liquidationPrice.String(), currentTwapPrice.String())
	}

	// Attempt Market Liquidation
	// This step will:
	// - Place market orders to close the perpetual.
	// - Execute Exchange -> SettleMarginAndRPnL -> OnPositionClose for each fill.
	// - OnPositionClose will handle the trader's PNL, margin refund, and use IF if netRefund for trader is negative (after accounting for fee).
	// - MarketLiquidation will flag for ADL if IF is insufficient OR if the order is not fully filled.
	closingRatio, err := k.MarketLiquidation(ctx, perpetual, market)
	if err != nil {
		// MarketLiquidation already flagged ADL if err was ErrInsufficientInsuranceFund.
		// Propagate the error; the caller (e.g., EndBlocker) might decide if ADL is the final step.
		return fmt.Errorf("forced_liquidation: MarketLiquidation failed for perp %d: %w", perpetual.Id, err)
	}

	//Reward on original margin amount.
	finalLiquidatorReward := market.LiquidationFeeShareRate.MulInt(perpetual.MarginAmount).TruncateInt()
	adlTriggeredByMarketLiquidation := false

	if closingRatio.IsZero() {
		// This means IF failed during MarketLiquidation's settlement process, ADL was set.
		// Policy: No liquidator fee if the core IF mechanism failed to backstop the position.
		// // flagged ADL, and returned (0, nil).
		finalLiquidatorReward = math.ZeroInt()
		adlTriggeredByMarketLiquidation = true
		ctx.Logger().Info(fmt.Sprintf("Liquidation for perpetual %d market %d failed due to insufficient IF, ADL set. No liquidator fee paid.", perpetual.Id, market.Id))
		// TODO: Emit event: LiquidationFailed_IF_Insolvent_ADL_Set
	} else if !closingRatio.Equal(math.LegacyOneDec()) && closingRatio.IsPositive() { // If partially filled and ratio > 0
		adlTriggeredByMarketLiquidation = true
		finalLiquidatorReward = closingRatio.MulInt(finalLiquidatorReward).TruncateInt()
		ctx.Logger().Info(fmt.Sprintf("ForcedLiquidation: Position %d (market %d) partially liquidated (ratio: %s), ADL set for remainder. Liquidator fee adjusted to %s.",
			perpetual.Id, market.Id, closingRatio.String(), finalLiquidatorReward.String()))
		// TODO: Emit event: LiquidationPartial_ADL_Set_For_Remainder_LiquidatorPaid
	}
	// 4. Pay Liquidator (if any reward is due)
	if finalLiquidatorReward.IsPositive() {
		// The Insurance Fund (via OnPositionClose) has already ensured the market.GetAccount()
		// is whole for the position's settlement, including covering the implicit cost of the fee.
		// Therefore, the market.GetAccount() is the correct source for this payment.
		errPay := k.bankKeeper.SendCoins(ctx, market.GetAccount(), liquidator, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, finalLiquidatorReward)))
		if errPay != nil {
			// This is a secondary failure. The position liquidation itself (or ADL flagging) has been processed.
			// Failure to pay the liquidator is critical but should not revert the primary liquidation resolution.
			// It might indicate an issue with the market account's general solvency or IF not fully covering.
			ctx.Logger().Error(fmt.Sprintf("CRITICAL: ForcedLiquidation: Failed to pay liquidator reward %s to %s from %s for perp %d after processing liquidation.",
				finalLiquidatorReward.String(), liquidator.String(), market.GetAccount().String(), perpetual.Id), "payment_error", errPay)
			// Ideally this should never happen as it has been covered in OnPositionClose
			// The main liquidation process is considered "handled" (either closed or ADL'd).
		} else {
			// TODO: Emit liquidator payment event
		}
	}

	if adlTriggeredByMarketLiquidation {
		// If MarketLiquidation returned (0, nil) due to IF error, it means ADL is the true resolution.
		// Return a specific error to signal this to the top-level caller.
		return types.ErrInsufficientInsuranceFund // Or a new error like ErrLiquidationRequiresADL
	}

	// TODO: Emit overarching LiquidationProcessed event (success, partial+ADL)
	// Note: Deletion of the fully liquidated perpetual is handled within Exchange -> SettleMarginAndRPnL.
	// If partially liquidated, the position quantity is updated, and ADL is set.
	return nil // Signifies the ForcedLiquidation attempt has been processed.
}

func (k Keeper) MarketLiquidation(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, error) {
	orderFilled := false
	closingRatio := math.LegacyOneDec()
	var err error

	msg := types.MsgPlaceMarketOrder{
		Creator:      perpetual.Owner,
		MarketId:     market.Id,
		BaseQuantity: perpetual.Quantity.Abs(),
		OrderType:    types.OrderType_ORDER_TYPE_MARKET_SELL,
	}
	// Even if equity value is 0 or -ve,
	// OnPositionClose internally handles by sending -ve amount from insurance fund to the market account
	cacheCtx, write := ctx.CacheContext()
	if perpetual.IsShort() {
		msg.OrderType = types.OrderType_ORDER_TYPE_MARKET_BUY
		orderFilled, err = k.ExecuteMarketBuyOrder(cacheCtx, market, msg, true)
	} else {
		orderFilled, err = k.ExecuteMarketSellOrder(cacheCtx, market, msg, true)
	}
	if err != nil {
		if errors.Is(err, types.ErrInsufficientInsuranceFund) {
			k.SetPerpetualADL(ctx, types.PerpetualADL{
				Id:       perpetual.Id,
				MarketId: perpetual.MarketId,
			})
		}
		return math.LegacyZeroDec(), nil
	}
	write()

	if !orderFilled {
		ctx.Logger().Warn(fmt.Sprintf("unable to fully liquidate %d for market %d, order cannot be filled", perpetual.Id, perpetual.MarketId))
		updatedPerpetual, err := k.GetPerpetual(ctx, perpetual.MarketId, perpetual.Id)
		if err != nil {
			return math.LegacyZeroDec(), err
		}
		closingRatio = (perpetual.Quantity.Abs().Sub(updatedPerpetual.Quantity.Abs())).Quo(perpetual.Quantity.Abs())

		// position is not fully closed
		k.SetPerpetualADL(ctx, types.PerpetualADL{
			Id:       perpetual.Id,
			MarketId: perpetual.MarketId,
		})
	}

	return closingRatio, nil
}
