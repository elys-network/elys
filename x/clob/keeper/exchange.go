package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Exchange(ctx sdk.Context, trade types.Trade) error {
	var buyerPerpetual types.Perpetual
	var err error
	buyerPerpetualOwner, buyerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.BuyerSubAccount.Id, trade.BuyerSubAccount.GetOwnerAccAddress(), trade.MarketId)
	if buyerAlreadyOwn {
		deleteBuyerPerpetual := false
		buyerPerpetual, err = k.GetPerpetual(ctx, trade.MarketId, buyerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}
		if buyerPerpetual.IsLong {
			buyerPerpetual.EntryPrice = buyerPerpetual.EntryPrice.Mul(buyerPerpetual.Quantity).Add(trade.Quantity.Mul(trade.Price)).Quo(buyerPerpetual.Quantity.Add(trade.Quantity))
			buyerPerpetual.Quantity = buyerPerpetual.Quantity.Add(trade.Quantity)
		} else {
			buyerPerpetual.Quantity = buyerPerpetual.Quantity.Sub(trade.Quantity)
			deleteBuyerPerpetual = buyerPerpetual.Quantity.IsZero()
			if buyerPerpetual.Quantity.IsNegative() {
				buyerPerpetual.Quantity = buyerPerpetual.Quantity.Abs()
				buyerPerpetual.IsLong = true
			}
		}
		if deleteBuyerPerpetual {
			k.DeletePerpetual(ctx, buyerPerpetual)
		} else {
			k.SetPerpetual(ctx, buyerPerpetual)
		}
	} else {
		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		buyerPerpetual = types.Perpetual{
			Id:           id,
			MarketId:     trade.MarketId,
			IsLong:       true,
			EntryPrice:   trade.Price,
			Owner:        trade.BuyerSubAccount.Owner,
			SubAccountId: trade.BuyerSubAccount.Id,
			Quantity:     trade.Quantity,
			Collateral:   trade.BuyerCollateral,
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

	var sellerPerpetual types.Perpetual

	sellerPerpetualOwner, sellerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.SellerSubAccount.Id, trade.SellerSubAccount.GetOwnerAccAddress(), trade.MarketId)
	if sellerAlreadyOwn {
		deleteSellerPerpetual := false
		sellerPerpetual, err = k.GetPerpetual(ctx, trade.MarketId, sellerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}
		if sellerPerpetual.IsLong {
			sellerPerpetual.EntryPrice = sellerPerpetual.EntryPrice.Mul(sellerPerpetual.Quantity).Add(trade.Quantity.Mul(trade.Price)).Quo(sellerPerpetual.Quantity.Add(trade.Quantity))
			sellerPerpetual.Quantity = sellerPerpetual.Quantity.Add(trade.Quantity)
		} else {
			sellerPerpetual.Quantity = sellerPerpetual.Quantity.Sub(trade.Quantity)
			deleteSellerPerpetual = sellerPerpetual.Quantity.IsZero()
			if sellerPerpetual.Quantity.IsNegative() {
				sellerPerpetual.Quantity = sellerPerpetual.Quantity.Abs()
				sellerPerpetual.IsLong = true
			}
		}
		if deleteSellerPerpetual {
			k.DeletePerpetual(ctx, sellerPerpetual)
		} else {
			k.SetPerpetual(ctx, sellerPerpetual)
		}
	} else {
		id := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		sellerPerpetual = types.Perpetual{
			Id:           id,
			MarketId:     trade.MarketId,
			IsLong:       true,
			EntryPrice:   trade.Price,
			Owner:        trade.SellerSubAccount.Owner,
			SubAccountId: trade.SellerSubAccount.Id,
			Quantity:     trade.Quantity,
			Collateral:   trade.SellerCollateral,
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
