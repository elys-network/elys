package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Exchange(ctx sdk.Context, trade types.Trade) error {
	buyerPerpetualOwner, buyerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.BuyerSubAccount.Id, trade.BuyerSubAccount.GetOwnerAccAddress(), trade.MarketId)
	if buyerAlreadyOwn {
		buyerPerpetual, err := k.GetPerpetual(ctx, trade.MarketId, buyerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}
		wasLong := buyerPerpetual.IsLong()
		oldEntryPrice := buyerPerpetual.EntryPrice
		oldQuantity := buyerPerpetual.Quantity
		buyerPerpetual.Quantity = buyerPerpetual.Quantity.Add(trade.Quantity)
		if buyerPerpetual.IsZero() {
			k.DeletePerpetual(ctx, buyerPerpetual)
		} else {
			if wasLong {
				buyerPerpetual.EntryPrice = oldEntryPrice.Mul(oldQuantity).Add(trade.Quantity.Mul(trade.Price)).Quo(buyerPerpetual.Quantity)
			} else {
				buyerPerpetual.EntryPrice = trade.Price
			}
			k.SetPerpetual(ctx, buyerPerpetual)
		}

	} else {
		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		buyerPerpetual := types.Perpetual{
			Id:           id,
			MarketId:     trade.MarketId,
			EntryPrice:   trade.Price,
			Owner:        trade.BuyerSubAccount.Owner,
			SubAccountId: trade.BuyerSubAccount.Id,
			Quantity:     trade.Quantity,
		}
		buyerPerpetualOwner = types.PerpetualOwner{
			Owner:        buyerPerpetual.Owner,
			SubAccountId: trade.BuyerSubAccount.Id,
			MarketId:     trade.MarketId,
			PerpetualId:  buyerPerpetual.Id,
		}
		k.SetPerpetual(ctx, buyerPerpetual)
		k.SetPerpetualOwner(ctx, buyerPerpetualOwner)
	}

	sellerPerpetualOwner, sellerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.SellerSubAccount.Id, trade.SellerSubAccount.GetOwnerAccAddress(), trade.MarketId)
	if sellerAlreadyOwn {
		sellerPerpetual, err := k.GetPerpetual(ctx, trade.MarketId, sellerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}
		wasShort := sellerPerpetual.IsShort()
		oldEntryPrice := sellerPerpetual.EntryPrice
		oldQuantity := sellerPerpetual.Quantity
		sellerPerpetual.Quantity = sellerPerpetual.Quantity.Sub(trade.Quantity)
		if sellerPerpetual.IsZero() {
			k.DeletePerpetual(ctx, sellerPerpetual)
		} else {
			if wasShort {
				sellerPerpetual.EntryPrice = oldEntryPrice.Mul(oldQuantity).Add(trade.Quantity.Neg().Mul(trade.Price)).Quo(sellerPerpetual.Quantity)
			} else {
				sellerPerpetual.EntryPrice = trade.Price
			}
			k.SetPerpetual(ctx, sellerPerpetual)
		}

	} else {
		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		sellerPerpetual := types.Perpetual{
			Id:           id,
			MarketId:     trade.MarketId,
			EntryPrice:   trade.Price,
			Owner:        trade.SellerSubAccount.Owner,
			SubAccountId: trade.SellerSubAccount.Id,
			Quantity:     trade.Quantity.Neg(),
		}
		sellerPerpetualOwner = types.PerpetualOwner{
			Owner:        sellerPerpetual.Owner,
			SubAccountId: trade.SellerSubAccount.Id,
			MarketId:     trade.MarketId,
			PerpetualId:  sellerPerpetual.Id,
		}
		k.SetPerpetual(ctx, sellerPerpetual)
		k.SetPerpetualOwner(ctx, sellerPerpetualOwner)
	}
	return nil
}
