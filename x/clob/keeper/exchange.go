package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Exchange(ctx sdk.Context, trade types.Trade) error {
	buyerPerpetual := types.Perpetual{}
	var err error
	if trade.BuyerPerpetualId != 0 {
		buyerPerpetual, err = k.GetPerpetual(ctx, trade.BuyerPerpetualId)
		if err != nil {
			return err
		}
		buyerPerpetual.EntryPrice = buyerPerpetual.EntryPrice.Mul(buyerPerpetual.Quantity).Add(trade.Quantity.Mul(trade.Price)).Quo(buyerPerpetual.Quantity.Add(trade.Quantity))
		buyerPerpetual.Quantity = buyerPerpetual.Quantity.Add(trade.Quantity)
	} else {
		buyerId := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		buyerPerpetual = types.Perpetual{
			Id:           buyerId,
			MarketId:     trade.MarketId,
			IsLong:       true,
			EntryPrice:   trade.Price,
			Owner:        trade.BuyerSubAccount.Owner,
			SubAccountId: trade.BuyerSubAccount.Id,
			Quantity:     trade.Quantity,
			Collateral:   trade.BuyerCollateral,
		}
	}
	sellerPerpetual := types.Perpetual{}
	if trade.SellerPerpetualId != 0 {
		sellerPerpetual, err = k.GetPerpetual(ctx, trade.SellerPerpetualId)
		if err != nil {
			return err
		}
		sellerPerpetual.EntryPrice = sellerPerpetual.EntryPrice.Mul(sellerPerpetual.Quantity).Add(trade.Quantity.Mul(trade.Price)).Quo(sellerPerpetual.Quantity.Add(trade.Quantity))
		sellerPerpetual.Quantity = sellerPerpetual.Quantity.Add(trade.Quantity)
	} else {
		sellerId := k.GetAndUpdatePerpetualCounter(ctx, trade.MarketId)
		sellerPerpetual = types.Perpetual{
			Id:           sellerId,
			MarketId:     trade.MarketId,
			IsLong:       false,
			EntryPrice:   trade.Price,
			Owner:        trade.SellerSubAccount.Owner,
			SubAccountId: trade.SellerSubAccount.Id,
			Quantity:     trade.Quantity,
			Collateral:   trade.SellerCollateral,
		}
	}
	k.SetPerpetual(ctx, buyerPerpetual)
	k.SetPerpetual(ctx, sellerPerpetual)
	return nil
}
