package keeper

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) ExecuteMarketBuyOrder(ctx sdk.Context, market types.PerpetualMarket, msg types.MsgPlaceMarketOrder, isLiquidation, isBuyerTaker bool) (bool, error) {
	if msg.OrderType != types.OrderType_ORDER_TYPE_MARKET_BUY {
		return false, errors.New("order is not a buy order")
	}
	var err error
	buyOrderFilled := false
	var buyerSubAccount types.SubAccount
	buyer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return false, err
	}

	sellIterator := k.GetSellOrderIterator(ctx, market.Id)
	buyerSubAccount, err = k.GetSubAccount(ctx, buyer, market.Id)
	if err != nil {
		return false, err
	}

	var sellOrdersToDelete []types.PerpetualOrderOwner
	filled := math.LegacyZeroDec()

	for ; sellIterator.Valid() && !buyOrderFilled; sellIterator.Next() {
		fmt.Println("---")
		var sellOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(sellIterator.Value(), &sellOrder)

		tradePrice := sellOrder.Price

		sellOrderFilled := false

		buyOrderMaxQuantity := msg.BaseQuantity.Sub(filled)
		sellOrderMaxQuantity := sellOrder.Amount.Sub(sellOrder.Filled)

		tradeQuantity := math.LegacyMinDec(buyOrderMaxQuantity, sellOrderMaxQuantity)
		if tradeQuantity.Equal(buyOrderMaxQuantity) {
			buyOrderFilled = true
		}
		if tradeQuantity.Equal(sellOrderMaxQuantity) {
			sellOrderFilled = true
		}

		sellOrder.Filled = sellOrder.Filled.Add(tradeQuantity)
		filled = filled.Add(tradeQuantity)

		if sellOrderFilled {
			toDelete := types.PerpetualOrderOwner{
				Owner:        sellOrder.Owner,
				SubAccountId: sellOrder.SubAccountId,
				OrderKey:     types.NewOrderKey(sellOrder.MarketId, sellOrder.OrderType, sellOrder.Price, sellOrder.BlockHeight),
			}
			sellOrdersToDelete = append(sellOrdersToDelete, toDelete)
		} else {
			k.SetPerpetualOrder(ctx, sellOrder)
		}
		fmt.Println("SELL ORDER EXECUTED: ")
		fmt.Println(sellOrder)

		sellerSubAccount, err := k.GetSubAccount(ctx, sellOrder.GetOwnerAccAddress(), market.Id)
		if err != nil {
			return false, err
		}
		err = k.Exchange(ctx, types.Trade{
			BuyerSubAccount:     buyerSubAccount,
			SellerSubAccount:    sellerSubAccount,
			MarketId:            market.Id,
			Price:               tradePrice,
			Quantity:            tradeQuantity,
			IsBuyerLiquidation:  isLiquidation,
			IsSellerLiquidation: false,
			IsBuyerTaker:        isBuyerTaker,
		})
		if err != nil {
			return false, err
		}
		fmt.Println("---")

	}

	err = sellIterator.Close()
	if err != nil {
		return false, err
	}
	for _, key := range sellOrdersToDelete {
		k.DeleteOrder(ctx, key.OrderKey)
		k.DeleteOrderOwner(ctx, types.PerpetualOrderOwner{
			Owner:        key.Owner,
			SubAccountId: key.SubAccountId,
			OrderKey:     key.OrderKey,
		})
	}

	return buyOrderFilled, nil

}

func (k Keeper) ExecuteMarketSellOrder(ctx sdk.Context, market types.PerpetualMarket, msg types.MsgPlaceMarketOrder, isLiquidation, isBuyerTaker bool) (bool, error) {
	if msg.OrderType != types.OrderType_ORDER_TYPE_MARKET_SELL {
		return false, errors.New("order is not a sell order")
	}
	var err error
	sellOrderFilled := false
	var sellerSubAccount types.SubAccount
	seller, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return false, err
	}

	buyIterator := k.GetBuyOrderIterator(ctx, market.Id)
	sellerSubAccount, err = k.GetSubAccount(ctx, seller, market.Id)
	if err != nil {
		return false, err
	}

	var buyOrdersToDelete []types.PerpetualOrderOwner
	filled := math.LegacyZeroDec()

	for ; buyIterator.Valid() && !sellOrderFilled; buyIterator.Next() {
		fmt.Println("---")
		var buyOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(buyIterator.Value(), &buyOrder)

		tradePrice := buyOrder.Price

		buyOrderFilled := false

		sellOrderMaxQuantity := msg.BaseQuantity.Sub(filled)
		buyOrderMaxQuantity := buyOrder.Amount.Sub(buyOrder.Filled)

		tradeQuantity := math.LegacyMinDec(buyOrderMaxQuantity, sellOrderMaxQuantity)
		if tradeQuantity.Equal(buyOrderMaxQuantity) {
			buyOrderFilled = true
		}
		if tradeQuantity.Equal(sellOrderMaxQuantity) {
			sellOrderFilled = true
		}

		buyOrder.Filled = buyOrder.Filled.Add(tradeQuantity)
		filled = filled.Add(tradeQuantity)

		if buyOrderFilled {
			toDelete := types.PerpetualOrderOwner{
				Owner:        buyOrder.Owner,
				SubAccountId: buyOrder.SubAccountId,
				OrderKey:     types.NewOrderKey(buyOrder.MarketId, buyOrder.OrderType, buyOrder.Price, buyOrder.BlockHeight),
			}
			buyOrdersToDelete = append(buyOrdersToDelete, toDelete)
		} else {
			k.SetPerpetualOrder(ctx, buyOrder)
		}
		fmt.Println("BUY ORDER EXECUTED: ")
		fmt.Println(buyOrder)

		buyerSubAccount, err := k.GetSubAccount(ctx, buyOrder.GetOwnerAccAddress(), market.Id)
		if err != nil {
			return false, err
		}
		err = k.Exchange(ctx, types.Trade{
			BuyerSubAccount:     buyerSubAccount,
			SellerSubAccount:    sellerSubAccount,
			MarketId:            market.Id,
			Price:               tradePrice,
			Quantity:            tradeQuantity,
			IsBuyerLiquidation:  false,
			IsSellerLiquidation: isLiquidation,
			IsBuyerTaker:        isBuyerTaker,
		})
		if err != nil {
			return false, err
		}
		fmt.Println("---")

	}

	err = buyIterator.Close()
	if err != nil {
		return false, err
	}
	for _, key := range buyOrdersToDelete {
		k.DeleteOrder(ctx, key.OrderKey)
		k.DeleteOrderOwner(ctx, types.PerpetualOrderOwner{
			Owner:        key.Owner,
			SubAccountId: key.SubAccountId,
			OrderKey:     key.OrderKey,
		})
	}

	return sellOrderFilled, nil

}
