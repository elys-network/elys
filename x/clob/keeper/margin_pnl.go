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
// 3. -ve to 0 (Buyer closing it's short position)
// 4. -ve to +ve (Buyer flipping sides short to long)
// 5. +ve to more +ve (Buyer expanding position)
// 6. +ve to less +ve (Seller reducing long position)
// 7. +ve to 0 (Seller closing it's long position)
// 8. +ve to -ve (Seller flipping sides long to short)
// 9. Opening of new position
// The function updates quantity, margin and entry price. Doesn't set in the KV store
// Funding rate needs to be updated outside before setting in KV
// Whenever there is AddToSubAccount, SettleRealizedPnL must be called as position is being reduced or closed
func (k Keeper) SettleMarginAndRPnL(ctx sdk.Context, market types.PerpetualMarket, oldPerpetual types.Perpetual, trade types.Trade, isBuyer bool) (updatedPerpetual types.Perpetual, err error) {

	if trade.Quantity.LTE(math.LegacyZeroDec()) {
		err = errors.New("trade quantity must be greater than zero")
		return
	}

	updatedPerpetual = oldPerpetual
	subAccount := trade.BuyerSubAccount
	updatedPerpetual.Quantity = oldPerpetual.Quantity.Add(trade.Quantity)
	if !isBuyer {
		updatedPerpetual.Quantity = oldPerpetual.Quantity.Sub(trade.Quantity)
		subAccount = trade.SellerSubAccount
	}

	// New Position is opened
	if oldPerpetual.Quantity.IsZero() || oldPerpetual.Quantity.IsNil() {
		requiredInitialMargin := trade.GetTradeValue().Mul(market.InitialMarginRatio).RoundInt()
		err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
		if err != nil {
			return
		}

		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		updatedPerpetual = types.Perpetual{
			Id:         id,
			MarketId:   trade.MarketId,
			EntryPrice: trade.Price,
			Owner:      subAccount.Owner,
			Quantity:   trade.Quantity,
			Margin:     requiredInitialMargin,
		}
		if !isBuyer {
			updatedPerpetual.Quantity = trade.Quantity.Neg()
		}

		return
	}

	// if its zero then it doesn't matter side.
	// This handles following:
	// 1.-ve to 0
	// 2.+ve to 0
	if updatedPerpetual.IsZero() {
		// Position fully closed, sending back short position margin
		err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, oldPerpetual.Margin)))
		if err != nil {
			return
		}

		// buyerPositionBefore is already negative
		err = k.SettleRealizedPnL(ctx, market, oldPerpetual.Quantity, subAccount, oldPerpetual.EntryPrice, trade.Price)
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
			newRequiredInitialMargin := updatedPerpetual.GetEntryValue().Mul(market.InitialMarginRatio).RoundInt()
			// TODO Edge case, Say IMR was 50%, previous margin was 500 (10 q, 100 p), now say IMR is changed to 1% and position expands to (11 q, 110 p)
			// TODO new margin required becomes 0.01*1210 = 12.1 which is less than old margin
			// Possible solution, refund (very risky, avoid), keep new margin same as old margin
			if newRequiredInitialMargin.LTE(oldPerpetual.Margin) {
				err = fmt.Errorf("newRequiredInitialMargin (%s) must be greater than oldPerpetual.Margin (%s) for buyer when position is increased from positive to more positive", newRequiredInitialMargin.String(), oldPerpetual.Margin.String())
				return
			}
			diff := newRequiredInitialMargin.Sub(oldPerpetual.Margin)
			err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
			if err != nil {
				return
			}

			// update margin
			updatedPerpetual.Margin = newRequiredInitialMargin
		} else {
			// Entry price remains same as position is reduced

			// Margin has to be reduced
			// We recalculate as IMR might have changed by gov
			newRequiredInitialMargin := updatedPerpetual.Quantity.Abs().Mul(updatedPerpetual.EntryPrice).Mul(market.InitialMarginRatio).RoundInt()
			// oldPerpetual.Margin must be greater than newRequiredInitialMargin
			if oldPerpetual.Margin.LTE(newRequiredInitialMargin) {
				err = fmt.Errorf("oldPerpetual.Margin (%s) must be greater than newRequiredInitialMargin (%s) for seller when position is reduced from positive to less positive", oldPerpetual.Margin.String(), newRequiredInitialMargin.String())
				return
			}
			diff := oldPerpetual.Margin.Sub(newRequiredInitialMargin)

			err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
			if err != nil {
				return
			}
			updatedPerpetual.Margin = newRequiredInitialMargin

			// position is partially closed for seller, 4 to 2, positionClosed = 2
			positionClosed := oldPerpetual.Quantity.Sub(updatedPerpetual.Quantity)
			err = k.SettleRealizedPnL(ctx, market, positionClosed, subAccount, oldPerpetual.EntryPrice, trade.Price)
			if err != nil {
				return
			}
		}
	}

	// definitely buyer side
	// This handles -ve to +ve
	if oldPerpetual.Quantity.IsNegative() && updatedPerpetual.Quantity.IsPositive() {
		// entry price becomes current trade price
		updatedPerpetual.EntryPrice = trade.Price

		// Return short position margin
		err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, oldPerpetual.Margin)))
		if err != nil {
			return
		}
		updatedPerpetual.Margin = math.ZeroInt()
		// Whole previous position gets closed, buyer old position is negative
		err = k.SettleRealizedPnL(ctx, market, oldPerpetual.Quantity, subAccount, oldPerpetual.EntryPrice, trade.Price)
		if err != nil {
			return
		}

		// Deduct long position margin
		requiredInitialMargin := updatedPerpetual.Quantity.Mul(trade.Price).Mul(market.InitialMarginRatio).RoundInt()
		err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
		if err != nil {
			return
		}
		updatedPerpetual.Margin = requiredInitialMargin
	}

	// This can happen on both side, buyer and seller
	// This handles following
	// 1. -ve to less -ve (buyer side)
	// 2. -ve to more -ve (seller side)
	if oldPerpetual.Quantity.IsNegative() && updatedPerpetual.Quantity.IsNegative() {
		if isBuyer {
			// Entry price remains same as position is reduced

			// Margin has to be reduced
			// We recalculate as IMR might have changed by gov
			newRequiredInitialMargin := updatedPerpetual.Quantity.Abs().Mul(updatedPerpetual.EntryPrice).Mul(market.InitialMarginRatio).RoundInt()
			// oldPerpetual.Margin must be greater than newRequiredInitialMargin
			if oldPerpetual.Margin.LTE(newRequiredInitialMargin) {
				err = fmt.Errorf("oldPerpetual.Margin (%s) must be greater than newRequiredInitialMargin (%s) for buyer when position is reduced from negative to less negative", oldPerpetual.Margin.String(), newRequiredInitialMargin.String())
				return
			}
			diff := oldPerpetual.Margin.Sub(newRequiredInitialMargin)

			err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
			if err != nil {
				return
			}
			updatedPerpetual.Margin = newRequiredInitialMargin

			// -4 to -2, positionClosed = -2
			positionClosed := oldPerpetual.Quantity.Sub(updatedPerpetual.Quantity)
			err = k.SettleRealizedPnL(ctx, market, positionClosed, subAccount, oldPerpetual.EntryPrice, trade.Price)
			if err != nil {
				return
			}

		} else {
			// entry price will become average
			num := oldPerpetual.GetEntryValue().Add(trade.GetTradeValue())
			updatedPerpetual.EntryPrice = num.Quo(updatedPerpetual.Quantity.Abs())

			// For margin, we will recalculate for the whole position as IMR can be changed, so just incremental won't work
			// Deducting margin
			// entry price MUST be updated before to average price
			newRequiredInitialMargin := updatedPerpetual.GetEntryValue().Mul(market.InitialMarginRatio).RoundInt()
			// TODO Edge case, Say IMR was 50%, previous margin was 500 (10 q, 100 p), now say IMR is changed to 1% and position expands to (11 q, 110 p)
			// TODO new margin required becomes 0.01*1210 = 12.1 which is less than old margin
			// Possible solution, refund (very risky, avoid), keep new margin same as old margin
			if newRequiredInitialMargin.LTE(oldPerpetual.Margin) {
				err = fmt.Errorf("newRequiredInitialMargin (%s) must be greater than oldPerpetual.Margin (%s) for seller when position is increased from negative to more negative", newRequiredInitialMargin.String(), oldPerpetual.Margin.String())
				return
			}
			diff := newRequiredInitialMargin.Sub(oldPerpetual.Margin)
			err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
			if err != nil {
				return
			}

			// update margin
			updatedPerpetual.Margin = newRequiredInitialMargin
		}
	}

	// definitely buyer side
	// This handles +ve to -ve
	if oldPerpetual.Quantity.IsPositive() && updatedPerpetual.Quantity.IsNegative() {
		// entry price becomes current trade price
		updatedPerpetual.EntryPrice = trade.Price

		// Seller reduces long position
		// Return long position margin
		err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, oldPerpetual.Margin)))
		if err != nil {
			return
		}
		updatedPerpetual.Margin = math.ZeroInt()

		// Whole previous position gets closed, seller old position is +ve
		err = k.SettleRealizedPnL(ctx, market, oldPerpetual.Quantity, subAccount, oldPerpetual.EntryPrice, trade.Price)
		if err != nil {
			return
		}

		// updatedPerpetual.Quantity is -ve
		requiredInitialMargin := updatedPerpetual.Quantity.Abs().Mul(trade.Price).Mul(market.InitialMarginRatio).RoundInt()
		err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
		if err != nil {
			return
		}
		updatedPerpetual.Margin = requiredInitialMargin

	}
	return
}

// SettleRealizedPnL settles the realized profit or loss for a closed position in a perpetual market.
// It transfers funds between the market and the subaccount based on the realized profit or loss.
// positionClosed should have sign, representing short or long
func (k Keeper) SettleRealizedPnL(ctx sdk.Context, market types.PerpetualMarket, positionClosed math.LegacyDec, subAccount types.SubAccount, entryPrice, tradePrice math.LegacyDec) (err error) {
	realizedPnl := positionClosed.Mul(tradePrice.Sub(entryPrice)).TruncateInt()
	if !realizedPnl.IsZero() {
		if realizedPnl.IsPositive() {
			// We send it from market because seller will send it to market as well
			err = k.AddToSubAccount(ctx, market.GetAccount(), subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, realizedPnl)))
			if err != nil {
				return err
			}
		} else {
			err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, realizedPnl.Neg())))
			if err != nil {
				return err
			}
		}
	}
	return
}
