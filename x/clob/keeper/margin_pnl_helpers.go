package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// handleNewPosition handles the opening of a new position
func (k Keeper) handleNewPosition(ctx sdk.Context, market types.PerpetualMarket, trade types.Trade,
	isBuyer bool, quoteDenomPrice math.LegacyDec) (types.Perpetual, error) {

	quantity := trade.Quantity
	subAccount := trade.BuyerSubAccount
	if !isBuyer {
		quantity = trade.Quantity.Neg()
		subAccount = trade.SellerSubAccount
	}

	// trade value x IMR / quote denom price
	requiredInitialMargin := trade.GetTradeValue().Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
	err := k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
	if err != nil {
		return types.Perpetual{}, err
	}

	return types.Perpetual{
		MarketId:     trade.MarketId,
		Owner:        subAccount.Owner,
		SubAccountId: subAccount.Id,
		Id:           k.GetAndIncrementPerpetualCounter(ctx, trade.MarketId),
		Quantity:     quantity,
		EntryPrice:   trade.Price,
		MarginAmount: requiredInitialMargin,
	}, nil
}

// handlePositionIncrease handles increasing an existing position
func (k Keeper) handlePositionIncrease(ctx sdk.Context, market types.PerpetualMarket,
	oldPerpetual types.Perpetual, trade types.Trade, quoteDenomPrice math.LegacyDec,
	subAccount types.SubAccount) (types.Perpetual, error) {

	updatedPerpetual := oldPerpetual
	updatedPerpetual.Quantity = updatedPerpetual.Quantity.Add(trade.Quantity)

	// Update entry price to weighted average
	totalValue := oldPerpetual.GetEntryValue().Add(trade.GetTradeValue())
	updatedPerpetual.EntryPrice = totalValue.Quo(updatedPerpetual.Quantity.Abs())

	// Calculate new required margin
	newRequiredInitialMargin := updatedPerpetual.GetEntryValue().Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()

	if newRequiredInitialMargin.LTE(oldPerpetual.MarginAmount) {
		return types.Perpetual{}, fmt.Errorf("new required margin %s must be greater than old margin %s when increasing position",
			newRequiredInitialMargin.String(), oldPerpetual.MarginAmount.String())
	}

	diff := newRequiredInitialMargin.Sub(oldPerpetual.MarginAmount)
	err := k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, diff)))
	if err != nil {
		return types.Perpetual{}, err
	}

	updatedPerpetual.MarginAmount = newRequiredInitialMargin
	return updatedPerpetual, nil
}

// handlePositionDecrease handles decreasing an existing position
func (k Keeper) handlePositionDecrease(ctx sdk.Context, market types.PerpetualMarket,
	oldPerpetual types.Perpetual, trade types.Trade, quoteDenomPrice math.LegacyDec,
	subAccount types.SubAccount, isLiquidation bool) (types.Perpetual, error) {

	updatedPerpetual := oldPerpetual

	if oldPerpetual.Quantity.IsPositive() {
		// Reducing long position (seller)
		updatedPerpetual.Quantity = updatedPerpetual.Quantity.Sub(trade.Quantity)
	} else {
		// Reducing short position (buyer)
		updatedPerpetual.Quantity = updatedPerpetual.Quantity.Add(trade.Quantity)
	}

	// Calculate PnL for the reduced portion
	var rpnl math.LegacyDec
	if oldPerpetual.Quantity.IsPositive() {
		rpnl = trade.Quantity.Mul(trade.Price.Sub(oldPerpetual.EntryPrice))
	} else {
		rpnl = trade.Quantity.Mul(oldPerpetual.EntryPrice.Sub(trade.Price))
	}

	// Calculate margin to return
	marginRatio := trade.Quantity.Quo(oldPerpetual.Quantity.Abs())
	marginToReturn := marginRatio.Mul(oldPerpetual.MarginAmount.ToLegacyDec())

	// Handle margin and PnL settlement
	err := k.settleMarginAndPnL(ctx, market, subAccount, marginToReturn, rpnl, quoteDenomPrice, isLiquidation)
	if err != nil {
		return types.Perpetual{}, err
	}

	// Update remaining margin
	updatedPerpetual.MarginAmount = oldPerpetual.MarginAmount.Sub(marginToReturn.RoundInt())

	return updatedPerpetual, nil
}

// handlePositionFlip handles when a position flips from long to short or vice versa
func (k Keeper) handlePositionFlip(ctx sdk.Context, market types.PerpetualMarket,
	oldPerpetual types.Perpetual, trade types.Trade, quoteDenomPrice math.LegacyDec,
	subAccount types.SubAccount, isLiquidation bool, newQuantity math.LegacyDec) (types.Perpetual, error) {

	// First close the old position
	err := k.OnPositionClose(ctx, market, oldPerpetual.Quantity, subAccount,
		oldPerpetual.EntryPrice, trade.Price, oldPerpetual.MarginAmount, quoteDenomPrice, isLiquidation)
	if err != nil {
		return types.Perpetual{}, err
	}

	// Create new position in opposite direction
	updatedPerpetual := oldPerpetual
	updatedPerpetual.Quantity = newQuantity
	updatedPerpetual.EntryPrice = trade.Price
	updatedPerpetual.MarginAmount = math.ZeroInt()

	// Calculate and collect margin for new position
	requiredInitialMargin := newQuantity.Abs().Mul(trade.Price).Mul(market.InitialMarginRatio).Quo(quoteDenomPrice).RoundInt()
	err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
	if err != nil {
		return types.Perpetual{}, err
	}

	updatedPerpetual.MarginAmount = requiredInitialMargin
	return updatedPerpetual, nil
}

// settleMarginAndPnL handles the settlement of margin and realized PnL
func (k Keeper) settleMarginAndPnL(ctx sdk.Context, market types.PerpetualMarket,
	subAccount types.SubAccount, marginToReturn math.LegacyDec, rpnl math.LegacyDec,
	quoteDenomPrice math.LegacyDec, isLiquidation bool) error {

	settlement := marginToReturn.Add(rpnl)
	settlementInt := settlement.Quo(quoteDenomPrice).RoundInt()

	if settlementInt.IsNegative() {
		// Market owes the trader - transfer from market to subaccount
		return k.bankKeeper.SendCoins(ctx, market.GetAccount(), subAccount.GetTradingAccountAddress(),
			sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, settlementInt.Abs())))
	} else if settlementInt.IsPositive() {
		// Trader owes the market
		return k.SendFromSubAccount(ctx, subAccount, market.GetAccount(),
			sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, settlementInt)))
	}

	return nil
}
