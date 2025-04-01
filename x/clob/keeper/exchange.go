package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
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

	requiredInitialMargin, err := trade.GetRequiredInitialMargin(market)
	if err != nil {
		return err
	}

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
		oldEntryPrice := buyerPerpetual.EntryPrice
		oldQuantity := buyerPerpetual.Quantity

		buyerPerpetual.Quantity = buyerPerpetual.Quantity.Add(trade.Quantity)
		if err != nil {
			return err
		}
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
				n1, err := oldEntryPrice.Mul(utils.IntToDec(oldQuantity))
				if err != nil {
					return err
				}
				n2, err := trade.QuantityDec().Mul(trade.Price)
				if err != nil {
					return err
				}
				num, err := n1.Add(n2)
				if err != nil {
					return err
				}
				buyerPerpetual.EntryPrice, err = num.Quo(buyerPerpetual.QunatityDec())
				if err != nil {
					return err
				}
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
			EntryFundingRate: market.CurrentFundingRate,
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
		oldEntryPrice := sellerPerpetual.EntryPrice
		oldQuantity := sellerPerpetual.Quantity
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
				n1, err := oldEntryPrice.Mul(utils.IntToDec(oldQuantity))
				if err != nil {
					return err
				}
				n2, err := utils.IntToDec(trade.Quantity.Neg()).Mul(trade.Price)
				if err != nil {
					return err
				}
				num, err := n1.Add(n2)
				if err != nil {
					return err
				}
				sellerPerpetual.EntryPrice, err = num.Quo(sellerPerpetual.QunatityDec())
				if err != nil {
					return err
				}
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
			EntryFundingRate: market.CurrentFundingRate,
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
	return nil
}
