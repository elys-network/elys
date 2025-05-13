package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

// SettleMarginAndRPnL Following cases possible:
// 1. -ve to more -ve (Seller expanding position)
// 2. -ve to less -ve (Buyer reducing position)
// 3. -ve to 0 (Buyer closing its short position)
// 4. -ve to +ve (Buyer flipping sides short to long)
// 5. +ve to more +ve (Buyer expanding position)
// 6. +ve to less +ve (Seller reducing long position)
// 7. +ve to 0 (Seller closing its long position)
// 8. +ve to -ve (Seller flipping sides-long to short)
// 9. Opening of new position
// The function updates quantity, margin and entry price. Doesn't set in the KV store.
// The Funding rate needs to be updated outside before setting in KV.
func (k Keeper) SettleMarginAndRPnL(ctx sdk.Context, market types.PerpetualMarket, oldPerpetual types.Perpetual, isLiquidation bool, trade types.Trade, isBuyer bool) (updatedPerpetual types.Perpetual, err error) {

	if trade.Quantity.LTE(math.LegacyZeroDec()) {
		err = errors.New("trade quantity must be greater than zero")
		return
	}

	quoteDenomPrice, err := k.GetDenomPrice(ctx, market.QuoteDenom)
	if err != nil {
		return
	}

	updatedPerpetual = oldPerpetual
	subAccount := trade.BuyerSubAccount

	// New Position is opened
	if oldPerpetual.Quantity.IsNil() || oldPerpetual.Quantity.IsZero() {
		quantity := trade.Quantity
		subAccount = trade.BuyerSubAccount
		if !isBuyer {
			quantity = trade.Quantity.Neg()
			subAccount = trade.SellerSubAccount
		}

		// trade value x IMR / quote denom price
		requiredInitialMargin := trade.GetTradeValue().Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
		err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
		if err != nil {
			return
		}

		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		updatedPerpetual = types.Perpetual{
			Id:           id,
			MarketId:     trade.MarketId,
			EntryPrice:   trade.Price,
			Owner:        subAccount.Owner,
			Quantity:     quantity,
			MarginAmount: requiredInitialMargin,
		}
		return
	}

	updatedPerpetual.Quantity = oldPerpetual.Quantity.Add(trade.Quantity)
	if !isBuyer {
		updatedPerpetual.Quantity = oldPerpetual.Quantity.Sub(trade.Quantity)
		subAccount = trade.SellerSubAccount
	}

	// If it's zero, then it doesn't matter side.
	// This handles the following:
	// 1.-ve to 0
	// 2.+ve to 0
	if updatedPerpetual.IsZero() {
		// Position fully closed, sending back short position margin
		// buyerPositionBefore is already negative
		err = k.OnPositionClose(ctx, market, oldPerpetual.Quantity, subAccount, oldPerpetual.EntryPrice, trade.Price, oldPerpetual.MarginAmount, quoteDenomPrice, isLiquidation)
		if err != nil {
			return
		}

		// position will get deleted
	}

	// This handles:
	// 1. +ve to more +ve (buyer)
	// 2. +ve to less +ve (seller)
	if oldPerpetual.Quantity.IsPositive() && updatedPerpetual.Quantity.IsPositive() {
		if isBuyer {
			// entry price will become average
			num := oldPerpetual.GetEntryValue().Add(trade.GetTradeValue())
			updatedPerpetual.EntryPrice = num.Quo(updatedPerpetual.Quantity)

			// For margin, we will recalculate for the whole position as IMR can be changed, so just incremental won't work
			// Deducting margin
			// entry price MUST be updated before to average price
			newRequiredInitialMargin := updatedPerpetual.GetEntryValue().Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
			// TODO Edge case, Say IMR was 50%, previous margin was 500 (10 q, 100 p), now say IMR is changed to 1% and position expands to (11 q, 110 p)
			// TODO new margin required becomes 0.01*1210 = 12.1 which is less than old margin
			// Possible solution, refund (very risky, avoid), keep new margin same as old margin
			if newRequiredInitialMargin.LTE(oldPerpetual.MarginAmount) {
				err = fmt.Errorf("newRequiredInitialMargin (%s) must be greater than oldPerpetual.Margin (%s) for buyer when position is increased from positive to more positive", newRequiredInitialMargin.String(), oldPerpetual.MarginAmount.String())
				return
			}
			diff := newRequiredInitialMargin.Sub(oldPerpetual.MarginAmount)
			err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
			if err != nil {
				return
			}

			// update margin
			updatedPerpetual.MarginAmount = newRequiredInitialMargin
		} else {
			// Entry price remains the same as position is reduced

			// The Margin has to be reduced
			// We recalculate as IMR might have changed by gov
			newRequiredInitialMargin := updatedPerpetual.Quantity.Abs().Mul(updatedPerpetual.EntryPrice).Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
			// oldPerpetual.Margin must be greater than newRequiredInitialMargin
			if oldPerpetual.MarginAmount.LTE(newRequiredInitialMargin) {
				err = fmt.Errorf("oldPerpetual.Margin (%s) must be greater than newRequiredInitialMargin (%s) for seller when position is reduced from positive to less positive", oldPerpetual.MarginAmount.String(), newRequiredInitialMargin.String())
				return
			}
			diff := oldPerpetual.MarginAmount.Sub(newRequiredInitialMargin)

			// position is partially closed for seller, 4 to 2, positionClosed = 2
			positionClosed := oldPerpetual.Quantity.Sub(updatedPerpetual.Quantity)

			err = k.OnPositionClose(ctx, market, positionClosed, subAccount, oldPerpetual.EntryPrice, trade.Price, diff, quoteDenomPrice, isLiquidation)
			if err != nil {
				return
			}

			updatedPerpetual.MarginAmount = newRequiredInitialMargin

		}
	}

	// definitely buyer side
	// This handles -ve to +ve
	if oldPerpetual.Quantity.IsNegative() && updatedPerpetual.Quantity.IsPositive() {
		// entry price becomes current trade price
		updatedPerpetual.EntryPrice = trade.Price

		// Return short position margin
		// The Whole previous position gets closed, buyer old position is negative
		err = k.OnPositionClose(ctx, market, oldPerpetual.Quantity, subAccount, oldPerpetual.EntryPrice, trade.Price, oldPerpetual.MarginAmount, quoteDenomPrice, isLiquidation)
		if err != nil {
			return
		}
		updatedPerpetual.MarginAmount = math.ZeroInt()

		// Deduct long position margin
		requiredInitialMargin := updatedPerpetual.Quantity.Mul(trade.Price).Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
		err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
		if err != nil {
			return
		}
		updatedPerpetual.MarginAmount = requiredInitialMargin
	}

	// This can happen on both sides, buyer and seller
	// This handles the following
	// 1. -ve to less -ve (buyer side)
	// 2. -ve to more -ve (seller side)
	if oldPerpetual.Quantity.IsNegative() && updatedPerpetual.Quantity.IsNegative() {
		if isBuyer {
			// Entry price remains the same as position is reduced

			// The Margin has to be reduced
			// We recalculate as IMR might have changed by gov
			newRequiredInitialMargin := updatedPerpetual.Quantity.Abs().Mul(updatedPerpetual.EntryPrice).Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
			// oldPerpetual.Margin must be greater than newRequiredInitialMargin
			if oldPerpetual.MarginAmount.LTE(newRequiredInitialMargin) {
				err = fmt.Errorf("oldPerpetual.Margin (%s) must be greater than newRequiredInitialMargin (%s) for buyer when position is reduced from negative to less negative", oldPerpetual.MarginAmount.String(), newRequiredInitialMargin.String())
				return
			}
			diff := oldPerpetual.MarginAmount.Sub(newRequiredInitialMargin)

			// -4 to -2, positionClosed = -2
			positionClosed := oldPerpetual.Quantity.Sub(updatedPerpetual.Quantity)

			err = k.OnPositionClose(ctx, market, positionClosed, subAccount, oldPerpetual.EntryPrice, trade.Price, diff, quoteDenomPrice, isLiquidation)
			if err != nil {
				return
			}

			updatedPerpetual.MarginAmount = newRequiredInitialMargin

		} else {
			// entry price will become average
			num := oldPerpetual.GetEntryValue().Add(trade.GetTradeValue())
			updatedPerpetual.EntryPrice = num.Quo(updatedPerpetual.Quantity.Abs())

			// For margin, we will recalculate for the whole position as IMR can be changed, so just incremental won't work
			// Deducting margin
			// entry price MUST be updated before to average price
			newRequiredInitialMargin := updatedPerpetual.GetEntryValue().Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
			// TODO Edge case, Say IMR was 50%, previous margin was 500 (10 q, 100 p), now say IMR is changed to 1% and position expands to (11 q, 110 p)
			// TODO new margin required becomes 0.01*1210 = 12.1 which is less than old margin
			// Possible solution, refund (very risky, avoid), keep new margin same as old margin
			if newRequiredInitialMargin.LTE(oldPerpetual.MarginAmount) {
				err = fmt.Errorf("newRequiredInitialMargin (%s) must be greater than oldPerpetual.Margin (%s) for seller when position is increased from negative to more negative", newRequiredInitialMargin.String(), oldPerpetual.MarginAmount.String())
				return
			}
			diff := newRequiredInitialMargin.Sub(oldPerpetual.MarginAmount)
			err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
			if err != nil {
				return
			}

			// update margin
			updatedPerpetual.MarginAmount = newRequiredInitialMargin
		}
	}

	// definitely buyer side
	// This handles +ve to -ve
	if oldPerpetual.Quantity.IsPositive() && updatedPerpetual.Quantity.IsNegative() {
		// entry price becomes current trade price
		updatedPerpetual.EntryPrice = trade.Price

		// Seller reduces long position
		// Return long position margin
		// Whole previous position gets closed, seller old position is +ve
		err = k.OnPositionClose(ctx, market, oldPerpetual.Quantity, subAccount, oldPerpetual.EntryPrice, trade.Price, oldPerpetual.MarginAmount, quoteDenomPrice, isLiquidation)
		if err != nil {
			return
		}

		updatedPerpetual.MarginAmount = math.ZeroInt()

		// updatedPerpetual.Quantity is -ve
		requiredInitialMargin := updatedPerpetual.Quantity.Abs().Mul(trade.Price).Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
		err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
		if err != nil {
			return
		}
		updatedPerpetual.MarginAmount = requiredInitialMargin

	}
	return
}

// OnPositionClose settles the realized profit or loss for a closed position in a perpetual market.
// It transfers funds from the market and to subAccount based on the realized profit or loss.
// positionClosed should have sign, representing short or long
func (k Keeper) OnPositionClose(ctx sdk.Context, market types.PerpetualMarket, positionClosed math.LegacyDec, subAccount types.SubAccount, entryPrice, tradePrice math.LegacyDec, refundMarginAmount math.Int, quoteDenomPrice math.LegacyDec, isLiquidation bool) (err error) {
	realizedPnl := positionClosed.Mul(tradePrice.Sub(entryPrice)).Quo(quoteDenomPrice).TruncateInt()
	netRefund := refundMarginAmount.Add(realizedPnl)

	if isLiquidation {
		liquidatorAmount := market.LiquidationFeeShareRate.MulInt(refundMarginAmount).TruncateInt()
		// if netRefund is already negative, this would make it more negative which will get compensated by the insurance fund and liquidator will get from market account
		netRefund = netRefund.Sub(liquidatorAmount)
	}

	if !netRefund.IsZero() {
		if netRefund.IsPositive() {
			// We send it from market because seller will send it to market as well
			err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, netRefund)))
			if err != nil {
				return err
			}
		} else {
			// This means losses exceed the collateral (initial margin)
			// and the position wasn't liquidated when the liquidation price was hit
			err = k.bankKeeper.SendCoins(ctx, market.GetInsuranceAccount(), market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, netRefund.Abs())))
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("failed to send amount(%s) from insurance fund to market (%d) fund on positon close", netRefund.Abs().String(), market.Id))
				// Don't think we can get error any other than insufficient funds
				return types.ErrInsufficientInsuranceFund
			}

			// TODO Trigger success event
		}
	} else {
		// TODO liquidation event
	}
	return
}
