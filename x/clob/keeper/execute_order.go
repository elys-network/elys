package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) ExecuteMarketBuyOrder(ctx sdk.Context, market types.PerpetualMarket, msg types.MsgPlaceMarketOrder, isLiquidation, isBuyerTaker bool) (bool, error) {
	if msg.OrderType != types.OrderType_ORDER_TYPE_MARKET_BUY {
		return false, types.ErrNotBuyOrder
	}
	var err error
	buyOrderFilled := false
	var buyerSubAccount types.SubAccount
	buyer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return false, err
	}

	sellIterator := k.GetSellOrderIterator(ctx, market.Id)
	subAccountId := types.CrossMarginSubAccountId
	if msg.IsIsolated {
		subAccountId = msg.MarketId
	}
	buyerSubAccount, err = k.GetSubAccount(ctx, buyer, subAccountId)
	if err != nil {
		return false, err
	}

	var sellOrdersToDelete []types.PerpetualOrderOwner
	filled := math.LegacyZeroDec()

	for ; sellIterator.Valid() && !buyOrderFilled; sellIterator.Next() {
		var sellOrder types.Order
		if err := k.cdc.Unmarshal(sellIterator.Value(), &sellOrder); err != nil {
			ctx.Logger().Error("failed to unmarshal sell order", "error", err)
			continue
		}

		tradePrice := sellOrder.GetPrice()

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
				OrderId:      types.NewOrderId(sellOrder.GetMarketId(), sellOrder.GetOrderType(), sellOrder.GetPriceTick(), sellOrder.GetCounter()),
			}
			sellOrdersToDelete = append(sellOrdersToDelete, toDelete)
		} else {
			err = k.SetOrder(ctx, sellOrder)
			if err != nil {
				return false, err
			}
		}

		sellerSubAccount, err := k.GetSubAccount(ctx, sellOrder.GetOwnerAccAddress(), sellOrder.SubAccountId)
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
	}

	err = sellIterator.Close()
	if err != nil {
		return false, err
	}
	// Batch delete all filled orders
	k.BatchDeleteOrders(ctx, sellOrdersToDelete)

	return buyOrderFilled, nil

}

func (k Keeper) ExecuteMarketSellOrder(ctx sdk.Context, market types.PerpetualMarket, msg types.MsgPlaceMarketOrder, isLiquidation, isBuyerTaker bool) (bool, error) {
	if msg.OrderType != types.OrderType_ORDER_TYPE_MARKET_SELL {
		return false, types.ErrNotSellOrder
	}
	var err error
	sellOrderFilled := false
	var sellerSubAccount types.SubAccount
	seller, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return false, err
	}

	subAccountId := types.CrossMarginSubAccountId
	if msg.IsIsolated {
		subAccountId = msg.MarketId
	}

	buyIterator := k.GetBuyOrderIterator(ctx, market.Id)
	sellerSubAccount, err = k.GetSubAccount(ctx, seller, subAccountId)
	if err != nil {
		return false, err
	}

	var buyOrdersToDelete []types.PerpetualOrderOwner
	filled := math.LegacyZeroDec()

	for ; buyIterator.Valid() && !sellOrderFilled; buyIterator.Next() {
		var buyOrder types.Order
		if err := k.cdc.Unmarshal(buyIterator.Value(), &buyOrder); err != nil {
			ctx.Logger().Error("failed to unmarshal buy order", "error", err)
			continue
		}

		tradePrice := buyOrder.GetPrice()

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
				OrderId:      types.NewOrderId(buyOrder.GetMarketId(), buyOrder.GetOrderType(), buyOrder.GetPriceTick(), buyOrder.GetCounter()),
			}
			buyOrdersToDelete = append(buyOrdersToDelete, toDelete)
		} else {
			err = k.SetOrder(ctx, buyOrder)
			if err != nil {
				return false, err
			}
		}

		buyerSubAccount, err := k.GetSubAccount(ctx, buyOrder.GetOwnerAccAddress(), buyOrder.SubAccountId)
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
	}

	err = buyIterator.Close()
	if err != nil {
		return false, err
	}
	// Batch delete all filled orders
	k.BatchDeleteOrders(ctx, buyOrdersToDelete)

	return sellOrderFilled, nil

}
