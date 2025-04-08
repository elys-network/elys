package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Exchange(ctx sdk.Context, trade types.Trade) error {

	newLong := false
	newShort := false
	longFullyClosed := false
	shortFullyClosed := false
	market, err := k.GetPerpetualMarket(ctx, trade.MarketId)
	if err != nil {
		return err
	}

	currentFundingRate := k.GetFundingRate(ctx, market.Id)

	requiredInitialMargin := trade.GetRequiredInitialMargin(market)

	buyerPerpetualOwner, buyerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.BuyerSubAccount.GetOwnerAccAddress(), trade.MarketId)
	// Deducting margin
	err = k.SendFromSubAccount(ctx, trade.BuyerSubAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
	if err != nil {
		return err
	}
	if buyerAlreadyOwn {
		buyerPerpetual, err := k.GetPerpetual(ctx, trade.MarketId, buyerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}

		err = k.SettleFunding(ctx, &trade.BuyerSubAccount, market, &buyerPerpetual)
		if err != nil {
			return err
		}

		wasLong := buyerPerpetual.IsLong()
		buyerOldEntryValue := buyerPerpetual.GetEntryValue()

		buyerPerpetual.Quantity = buyerPerpetual.Quantity.Add(trade.Quantity)
		if buyerPerpetual.IsZero() {
			// Buyer was in short position
			k.DeletePerpetual(ctx, buyerPerpetual)
			// Sending back short position margin
			err = k.AddToSubAccount(ctx, market.GetAccount(), trade.BuyerSubAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, buyerPerpetual.Margin)))
			if err != nil {
				return err
			}
			longFullyClosed = true
		} else {
			if wasLong {
				num := buyerOldEntryValue.Add(trade.GetTradeValue())
				buyerPerpetual.EntryPrice = num.Quo(buyerPerpetual.Quantity.ToLegacyDec())
			} else {
				buyerPerpetual.EntryPrice = trade.Price
			}
			k.SetPerpetual(ctx, buyerPerpetual)
		}

	} else {
		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		buyerPerpetual := types.Perpetual{
			Id:               id,
			MarketId:         trade.MarketId,
			EntryPrice:       trade.Price,
			Owner:            trade.BuyerSubAccount.Owner,
			Quantity:         trade.Quantity,
			Margin:           requiredInitialMargin,
			EntryFundingRate: currentFundingRate.Rate,
		}
		buyerPerpetualOwner = types.PerpetualOwner{
			Owner:       buyerPerpetual.Owner,
			MarketId:    trade.MarketId,
			PerpetualId: buyerPerpetual.Id,
		}
		k.SetPerpetual(ctx, buyerPerpetual)
		k.SetPerpetualOwner(ctx, buyerPerpetualOwner)
		newLong = true
	}

	sellerPerpetualOwner, sellerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.SellerSubAccount.GetOwnerAccAddress(), trade.MarketId)
	err = k.SendFromSubAccount(ctx, trade.SellerSubAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, requiredInitialMargin)))
	if err != nil {
		return err
	}
	if sellerAlreadyOwn {

		sellerPerpetual, err := k.GetPerpetual(ctx, trade.MarketId, sellerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}

		err = k.SettleFunding(ctx, &trade.SellerSubAccount, market, &sellerPerpetual)
		if err != nil {
			return err
		}

		wasShort := sellerPerpetual.IsShort()
		sellerOldEntryValue := sellerPerpetual.GetEntryValue()
		sellerPerpetual.Quantity = sellerPerpetual.Quantity.Sub(trade.Quantity)
		if sellerPerpetual.IsZero() {
			k.DeletePerpetual(ctx, sellerPerpetual)
			err = k.AddToSubAccount(ctx, market.GetAccount(), trade.BuyerSubAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, sellerPerpetual.Margin)))
			if err != nil {
				return err
			}
			shortFullyClosed = true
		} else {
			if wasShort {
				num := sellerOldEntryValue.Add(trade.GetTradeValue())
				sellerPerpetual.EntryPrice = num.Quo(sellerPerpetual.Quantity.ToLegacyDec())
			} else {
				sellerPerpetual.EntryPrice = trade.Price
			}
			k.SetPerpetual(ctx, sellerPerpetual)
		}

	} else {
		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		sellerPerpetual := types.Perpetual{
			Id:               id,
			MarketId:         trade.MarketId,
			EntryPrice:       trade.Price,
			Owner:            trade.SellerSubAccount.Owner,
			Quantity:         trade.Quantity.Neg(),
			Margin:           requiredInitialMargin,
			EntryFundingRate: currentFundingRate.Rate,
		}
		sellerPerpetualOwner = types.PerpetualOwner{
			Owner:       sellerPerpetual.Owner,
			MarketId:    trade.MarketId,
			PerpetualId: sellerPerpetual.Id,
		}
		k.SetPerpetual(ctx, sellerPerpetual)
		k.SetPerpetualOwner(ctx, sellerPerpetualOwner)
		newShort = true
	}

	if newLong && newShort {
		market.TotalOpen = market.TotalOpen.Add(trade.Quantity)
	}
	if longFullyClosed && shortFullyClosed {
		market.TotalOpen = market.TotalOpen.Sub(trade.Quantity)
	}
	k.SetPerpetualMarket(ctx, market)
	k.SetTwapPrices(ctx, trade)
	return nil
}
