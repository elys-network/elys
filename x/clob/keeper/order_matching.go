package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) ExecuteMarket(ctx sdk.Context, marketId uint64) error {
	market, err := k.GetPerpetualMarket(ctx, marketId)
	if err != nil {
		return err
	}

	fullyFilled := true
	buyOrderIterator := k.GetBuyOrderIterator(ctx, marketId)

	var keysToDelete [][]byte

	for ; buyOrderIterator.Valid() && fullyFilled; buyOrderIterator.Next() {
		var buyOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(buyOrderIterator.Value(), &buyOrder)
		fullyFilled, err = k.ExecuteLimitBuyOrder(ctx, market, &buyOrder)
		if err != nil {
			return err
		}
		if fullyFilled {
			// iterator.Key() gives key bytes without prefix
			keysToDelete = append(keysToDelete, types.GetPerpetualOrderKey(buyOrder.MarketId, buyOrder.OrderType, buyOrder.Price, buyOrder.BlockHeight))
		}
	}

	err = buyOrderIterator.Close()
	if err != nil {
		return err
	}
	for _, key := range keysToDelete {
		k.DeleteOrder(ctx, key)
	}

	return nil
}

func (k Keeper) ExecuteLimitBuyOrder(ctx sdk.Context, market types.PerpetualMarket, buyOrder *types.PerpetualOrder) (bool, error) {
	if !buyOrder.IsBuy() {
		return false, errors.New("order is not a buy order")
	}
	var err error
	buyOrderFilled := false
	var buyerSubAccount, sellerSubAccount types.SubAccount

	highestBuyPrice := buyOrder.Price
	//lowestSellPrice := k.GetLowestSellPrice(ctx, market.Id)

	sellIterator := k.GetSellOrderIterator(ctx, market.Id)
	buyerSubAccount, err = k.GetSubAccount(ctx, buyOrder.GetOwnerAccAddress(), market.Id)
	if err != nil {
		return false, err
	}

	var keysToDelete [][]byte

	stop := false

	for ; sellIterator.Valid() && !buyOrderFilled && !stop; sellIterator.Next() {
		fmt.Println("---")
		var sellOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(sellIterator.Value(), &sellOrder)

		if sellOrder.IsBuy() {
			continue
		}

		lowestSellPrice := sellOrder.Price

		if highestBuyPrice.GTE(lowestSellPrice) {
			sellOrderFilled := false

			tradePrice := sellOrder.Price
			if sellOrder.BlockHeight > buyOrder.BlockHeight {
				tradePrice = buyOrder.Price
			}
			if sellOrder.BlockHeight == buyOrder.BlockHeight {
				sumPrice := buyOrder.Price.Add(sellOrder.Price)
				tradePrice = sumPrice.Quo(math.LegacyNewDec(2))
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
				// iterator.Key() gives key bytes without prefix
				keysToDelete = append(keysToDelete, types.GetPerpetualOrderKey(sellOrder.MarketId, sellOrder.OrderType, sellOrder.Price, sellOrder.BlockHeight))
			} else {
				k.SetPerpetualOrder(ctx, sellOrder)
			}
			fmt.Println("SELL ORDER EXECUTED: ")
			fmt.Println(sellOrder)

			sellerSubAccount, err = k.GetSubAccount(ctx, sellOrder.GetOwnerAccAddress(), market.Id)
			if err != nil {
				return false, err
			}
			err = k.Exchange(ctx, types.Trade{
				BuyerSubAccount:  buyerSubAccount,
				SellerSubAccount: sellerSubAccount,
				MarketId:         market.Id,
				Price:            tradePrice,
				Quantity:         tradeQuantity,
				IsLiquidation:    false,
			})
			if err != nil {
				return false, err
			}
			fmt.Println("BUY ORDER EXECUTED: ")
			fmt.Println(*buyOrder)
			fmt.Println("---")
		} else {
			stop = true
		}

	}

	err = sellIterator.Close()
	if err != nil {
		return false, err
	}
	for _, key := range keysToDelete {
		k.DeleteOrder(ctx, key)
	}

	if !buyOrderFilled {
		k.SetPerpetualOrder(ctx, *buyOrder)
	}

	return buyOrderFilled, nil

}
