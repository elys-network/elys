package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

// ForcedLiquidation Possible cases in isolated margin:
//  1. Equity Value > 0 and the market absorbs the whole position - liquidator fee paid by margin
//  2. Equity Value > 0 and the market only partially absorbs position - use ADL to close the rest of the position and partial liquidator fee can be paid by margin
//  3. Equity Value <= 0 and there are enough insurance funds, the market absorbs the whole position - liquidator fee paid by insurance funds
//  4. Equity Value <= 0 and there are enough insurance funds, the market partially absorbs the whole position - liquidator fee paid by insurance funds
//  5. Equity Value <= 0 and there are NOT enough insurance funds.
//
// Auto Deleverage: Find positions with maximum profit and close against that.
// Steps:
// - Place market orders to close the perpetual.
// - Execute Exchange -> SettleMarginAndRPnL -> OnPositionClose for each fill.
// - OnPositionClose will handle the trader's PNL, margin refund, and use IF if netRefund for trader is negative (after accounting for fee).
// - MarketLiquidation will flag for ADL if IF is insufficient OR if the order is not fully filled.
func (k Keeper) ForcedLiquidation(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket, liquidator sdk.AccAddress) (math.Int, error) {
	subAccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), perpetual.SubAccountId)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("forced_liquidation: failed to get subaccount for owner %s, market %d: %w", perpetual.Owner, market.Id, err)
	}
	liquidationPrice, err := k.GetLiquidationPrice(ctx, perpetual, market, subAccount)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("forced_liquidation: failed to get liquidation price for perp %d: %w", perpetual.Id, err)
	}
	currentTwapPrice := k.GetCurrentTwapPrice(ctx, market.Id) // Assuming this is the Mark Price
	if currentTwapPrice.IsZero() {
		return math.ZeroInt(), errors.New("forced_liquidation: cannot liquidate, current twap mark price is zero")
	}
	if perpetual.IsLong() && liquidationPrice.LT(currentTwapPrice) { // Price hasn't dropped enough for Long
		return math.ZeroInt(), fmt.Errorf("forced_liquidation: long position not yet liquidatable. LiqPx: %s, MarkPx: %s", liquidationPrice.String(), currentTwapPrice.String())
	}
	if perpetual.IsShort() && liquidationPrice.GT(currentTwapPrice) { // Price hasn't risen enough for Short
		return math.ZeroInt(), fmt.Errorf("forced_liquidation: short position not yet liquidatable. LiqPx: %s, MarkPx: %s", liquidationPrice.String(), currentTwapPrice.String())
	}

	closingRatio, adlTriggered, err := k.MarketLiquidation(ctx, perpetual, market)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("forced_liquidation: MarketLiquidation failed for perp %d: %w", perpetual.Id, err)
	}

	// This indicates that either the IF was not enough during the settlement of the filled portion (closing ratio 0),
	// or the position was only partially filled and the remainder needs ADL (0 < closing ratio < 1).
	if adlTriggered {
		k.SetPerpetualADL(ctx, types.PerpetualADL{
			Id:       perpetual.Id,
			MarketId: perpetual.MarketId,
		})

		ctx.Logger().Info(fmt.Sprintf("ForcedLiquidation: Position %d (market %d) (liquidation closing ratio: %s), ADL set", perpetual.Id, market.Id, closingRatio.String()))

		// TODO: Emit event: LiquidationRequiresADL with closing ratio
	} else {
		// A case where say a previous liquidation triggers ADL due to lack of insurance fund or lack of orders, but now liquidtion can be done again.
		k.DeletePerpetualADL(ctx, perpetual.MarketId, perpetual.Id)
	}

	if closingRatio.IsZero() {
		// This means IF failed during MarketLiquidation's settlement process, ADL was set.
		// No liquidator fee if the core IF mechanism failed to backstop the position.
		// We want to return nil here because we want to commit adl for the tx

		// TODO: Emit overarching LiquidationProcessed event (success, partial+ADL) with closing ratio
		return math.ZeroInt(), nil
	}

	// LiquidationFeeShareRate * MarginAmount * closingRatio
	finalLiquidatorReward := market.LiquidationFeeShareRate.MulInt(perpetual.MarginAmount).Mul(closingRatio).TruncateInt()

	err = k.bankKeeper.SendCoins(ctx, market.GetAccount(), liquidator, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, finalLiquidatorReward)))
	if err != nil {
		// Ideally, this should never happen as it has been covered in OnPositionClose. If liquidator payment fails, something fundamental is wrong
		ctx.Logger().Error(fmt.Sprintf("CRITICAL: ForcedLiquidation: Failed to pay liquidator reward %s to %s from %s for perp %d after processing liquidation. payment_error: %s", finalLiquidatorReward.String(), liquidator.String(), market.GetAccount().String(), perpetual.Id, err.Error()))
		return finalLiquidatorReward, err
	}

	// TODO: Emit overarching LiquidationProcessed event (success, partial+ADL) with closing ratio
	return finalLiquidatorReward, nil
}

func (k Keeper) MarketLiquidation(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, bool, error) {
	orderFilled := false
	closingRatio := math.LegacyOneDec()
	// There can be a case where Market execution fails not due to lack of funds or lack of orders
	adlTriggered := false
	var err error

	msg := types.MsgPlaceMarketOrder{
		Creator:      perpetual.Owner,
		MarketId:     market.Id,
		BaseQuantity: perpetual.Quantity.Abs(),
		OrderType:    types.OrderType_ORDER_TYPE_MARKET_SELL,
	}
	// Even if equity value is 0 or -ve, OnPositionClose internally handles by sending -ve amount from insurance fund to the market account
	cacheCtx, write := ctx.CacheContext()
	if perpetual.IsShort() {
		msg.OrderType = types.OrderType_ORDER_TYPE_MARKET_BUY
		orderFilled, err = k.ExecuteMarketBuyOrder(cacheCtx, market, msg, true, false)
	} else {
		orderFilled, err = k.ExecuteMarketSellOrder(cacheCtx, market, msg, true, false)
	}
	if err != nil {
		if errors.Is(err, types.ErrInsufficientInsuranceFund) {
			err = nil
			adlTriggered = true
		}
		// Nothing got liquidated, return 0, nil.
		// We return nil error here because we want to commit adl for the tx
		return math.LegacyZeroDec(), adlTriggered, err
	}
	write()

	if !orderFilled {
		ctx.Logger().Warn(fmt.Sprintf("unable to fully liquidate %d for market %d, order cannot be filled", perpetual.Id, perpetual.MarketId))
		updatedPerpetual, err := k.GetPerpetual(ctx, perpetual.MarketId, perpetual.Id)
		if err != nil {
			return math.LegacyZeroDec(), false, err
		}
		closingRatio = (perpetual.Quantity.Abs().Sub(updatedPerpetual.Quantity.Abs())).Quo(perpetual.Quantity.Abs())
		adlTriggered = true
	}

	return closingRatio, adlTriggered, nil
}
