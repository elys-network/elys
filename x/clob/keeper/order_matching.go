package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) ExecuteMarket(ctx sdk.Context, marketId uint64) error {
	market, err := k.GetPerpetualMarket(ctx, marketId)
	if err != nil {
		return err
	}

	fullyFilled := true
	buyOrderIterator := k.GetBuyOrderIterator(ctx, marketId)

	var keysToDelete []types.PerpetualOrderOwner

	for ; buyOrderIterator.Valid() && fullyFilled; buyOrderIterator.Next() {
		var buyOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(buyOrderIterator.Value(), &buyOrder)
		fullyFilled, err = k.ExecuteLimitBuyOrder(ctx, market, &buyOrder)
		if err != nil {
			return err
		}
		if fullyFilled {
			toDelete := types.PerpetualOrderOwner{
				Owner:        buyOrder.Owner,
				SubAccountId: buyOrder.SubAccountId,
				OrderKey:     types.NewOrderKey(buyOrder.MarketId, buyOrder.OrderType, buyOrder.Price, buyOrder.Counter),
			}
			keysToDelete = append(keysToDelete, toDelete)
		}
	}

	err = buyOrderIterator.Close()
	if err != nil {
		return err
	}
	// Batch delete all filled orders
	k.BatchDeleteOrders(ctx, keysToDelete)

	return nil
}

func (k Keeper) ExecuteLimitBuyOrder(ctx sdk.Context, market types.PerpetualMarket, buyOrder *types.PerpetualOrder) (bool, error) {
	if !buyOrder.IsBuy() {
		return false, types.ErrNotBuyOrder
	}
	var err error
	buyOrderFilled := false
	var buyerSubAccount, sellerSubAccount types.SubAccount

	highestBuyPrice := buyOrder.Price
	//lowestSellPrice := k.GetLowestSellPrice(ctx, market.Id)

	sellIterator := k.GetSellOrderIterator(ctx, market.Id)
	buyerSubAccount, err = k.GetSubAccount(ctx, buyOrder.GetOwnerAccAddress(), buyOrder.SubAccountId)
	if err != nil {
		return false, err
	}

	var keysToDelete []types.PerpetualOrderOwner

	stop := false

	for ; sellIterator.Valid() && !buyOrderFilled && !stop; sellIterator.Next() {
		var sellOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(sellIterator.Value(), &sellOrder)

		if sellOrder.IsBuy() {
			continue
		}

		lowestSellPrice := sellOrder.Price

		if highestBuyPrice.GTE(lowestSellPrice) {
			sellOrderFilled := false

			tradePrice := sellOrder.Price
			if sellOrder.Counter > buyOrder.Counter {
				tradePrice = buyOrder.Price
			}

			// remainingQuantity = buyOrderQuantity at trade price - already filled
			buyOrderMaxQuantity := buyOrder.Amount.Sub(buyOrder.Filled)
			sellOrderMaxQuantity := sellOrder.Amount.Sub(sellOrder.Filled)

			tradeQuantity := math.LegacyMinDec(buyOrderMaxQuantity, sellOrderMaxQuantity)
			if tradeQuantity.Equal(buyOrderMaxQuantity) {
				buyOrderFilled = true
			}
			if tradeQuantity.Equal(sellOrderMaxQuantity) {
				sellOrderFilled = true
			}

			sellOrder.Filled = sellOrder.Filled.Add(tradeQuantity)
			buyOrder.Filled = buyOrder.Filled.Add(tradeQuantity)

			if sellOrderFilled {
				toDelete := types.PerpetualOrderOwner{
					Owner:        sellOrder.Owner,
					SubAccountId: sellOrder.SubAccountId,
					OrderKey:     types.NewOrderKey(sellOrder.MarketId, sellOrder.OrderType, sellOrder.Price, sellOrder.Counter),
				}
				keysToDelete = append(keysToDelete, toDelete)
			} else {
				k.SetPerpetualOrder(ctx, sellOrder)
			}

			sellerSubAccount, err = k.GetSubAccount(ctx, sellOrder.GetOwnerAccAddress(), sellOrder.SubAccountId)
			if err != nil {
				return false, err
			}

			isBuyerTaker := buyOrder.Counter > sellOrder.Counter
			err = k.Exchange(ctx, types.Trade{
				BuyerSubAccount:     buyerSubAccount,
				SellerSubAccount:    sellerSubAccount,
				MarketId:            market.Id,
				Price:               tradePrice,
				Quantity:            tradeQuantity,
				IsBuyerLiquidation:  false,
				IsSellerLiquidation: false,
				IsBuyerTaker:        isBuyerTaker,
			})
			if err != nil {
				return false, err
			}
		} else {
			stop = true
		}

	}

	err = sellIterator.Close()
	if err != nil {
		return false, err
	}
	// Batch delete all filled orders
	k.BatchDeleteOrders(ctx, keysToDelete)

	if !buyOrderFilled {
		k.SetPerpetualOrder(ctx, *buyOrder)
	}

	return buyOrderFilled, nil

}
