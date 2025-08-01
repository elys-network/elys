package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) ExecuteMarket(ctx sdk.Context, marketId uint64, limit int64) (totalVolumeValue math.LegacyDec, ordersProcessed uint64, totalVolume math.LegacyDec, err error) {
	totalVolume = math.LegacyZeroDec()
	totalVolumeValue = math.LegacyZeroDec()

	market, err := k.GetPerpetualMarket(ctx, marketId)
	if err != nil {
		return totalVolumeValue, ordersProcessed, totalVolume, err
	}

	fullyFilled := true
	buyOrderIterator := k.GetBuyOrderIterator(ctx, marketId)

	var keysToDelete []types.PerpetualOrderOwner

	totalSellOrdersExecuted := int64(0)

	for ; buyOrderIterator.Valid() && fullyFilled; buyOrderIterator.Next() {
		if limit <= 0 {
			break
		}
		var buyOrder types.PerpetualOrder
		k.cdc.MustUnmarshal(buyOrderIterator.Value(), &buyOrder)

		sellOrdersExecuted := int64(0)
		currentTradeVolume := math.LegacyZeroDec()
		currentTradeValue := math.LegacyZeroDec()
		currentTradeValue, sellOrdersExecuted, currentTradeVolume, fullyFilled, err = k.ExecuteLimitBuyOrder(ctx, market, &buyOrder)
		if err != nil {
			return totalVolumeValue, ordersProcessed, totalVolume, err
		}

		totalSellOrdersExecuted = totalSellOrdersExecuted + sellOrdersExecuted

		limit = limit - sellOrdersExecuted

		if currentTradeVolume.IsPositive() {
			totalVolume = totalVolume.Add(currentTradeVolume)
			totalVolumeValue = totalVolumeValue.Add(currentTradeValue)
			ordersProcessed++
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
		return totalVolumeValue, ordersProcessed, totalVolume, err
	}
	// Batch delete all filled orders
	k.BatchDeleteOrders(ctx, keysToDelete)

	// Emit event for market execution
	if ordersProcessed > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventOrderBookUpdate,
				sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", marketId)),
				sdk.NewAttribute("buy_orders_processed", fmt.Sprintf("%d", ordersProcessed)),
				sdk.NewAttribute("total_volume", totalVolume.String()),
				sdk.NewAttribute("total_volume_value", totalVolumeValue.String()),
				sdk.NewAttribute("orders_deleted", fmt.Sprintf("%d", len(keysToDelete))),
			),
		)
	}

	return totalVolumeValue, ordersProcessed, totalVolume, nil
}

func (k Keeper) ExecuteLimitBuyOrder(ctx sdk.Context, market types.PerpetualMarket, buyOrder *types.PerpetualOrder) (totalVolumeValue math.LegacyDec, tradesExecuted int64, totalTradeVolume math.LegacyDec, buyOrderFilled bool, err error) {
	totalTradeVolume = math.LegacyZeroDec()
	totalVolumeValue = math.LegacyZeroDec()
	buyOrderFilled = false

	if !buyOrder.IsBuy() {
		err = types.ErrNotBuyOrder
		return
	}
	var buyerSubAccount, sellerSubAccount types.SubAccount

	highestBuyPrice := buyOrder.Price
	//lowestSellPrice := k.GetLowestSellPrice(ctx, market.Id)

	sellIterator := k.GetSellOrderIterator(ctx, market.Id)
	buyerSubAccount, err = k.GetSubAccount(ctx, buyOrder.GetOwnerAccAddress(), buyOrder.SubAccountId)
	if err != nil {
		return
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
				err = k.SetPerpetualOrder(ctx, sellOrder)
				if err != nil {
					return
				}
			}

			sellerSubAccount, err = k.GetSubAccount(ctx, sellOrder.GetOwnerAccAddress(), sellOrder.SubAccountId)
			if err != nil {
				return
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
				return
			}

			// Track trade metrics
			tradesExecuted++
			totalTradeVolume = totalTradeVolume.Add(tradeQuantity)
			totalVolumeValue = totalVolumeValue.Add(tradePrice.Mul(tradeQuantity))

			// Emit trade event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTrade,
					sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", market.Id)),
					sdk.NewAttribute(types.AttributeBuyer, buyOrder.Owner),
					sdk.NewAttribute(types.AttributeSeller, sellOrder.Owner),
					sdk.NewAttribute(types.AttributeTradePrice, tradePrice.String()),
					sdk.NewAttribute(types.AttributeTradeQuantity, tradeQuantity.String()),
					sdk.NewAttribute(types.AttributeIsTaker, fmt.Sprintf("%t", isBuyerTaker)),
					sdk.NewAttribute("buy_order_counter", fmt.Sprintf("%d", buyOrder.Counter)),
					sdk.NewAttribute("sell_order_counter", fmt.Sprintf("%d", sellOrder.Counter)),
				),
			)
		} else {
			stop = true
		}

	}

	err = sellIterator.Close()
	if err != nil {
		return
	}
	// Batch delete all filled orders
	k.BatchDeleteOrders(ctx, keysToDelete)

	if !buyOrderFilled {
		err = k.SetPerpetualOrder(ctx, *buyOrder)
		if err != nil {
			return
		}
	}

	// Emit order execution event if trades were made
	if tradesExecuted > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventOrderExecuted,
				sdk.NewAttribute(types.AttributeOrderId, fmt.Sprintf("%d", buyOrder.Counter)),
				sdk.NewAttribute(types.AttributeOwner, buyOrder.Owner),
				sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", buyOrder.MarketId)),
				sdk.NewAttribute(types.AttributeOrderType, buyOrder.OrderType.String()),
				sdk.NewAttribute(types.AttributeFilledQuantity, buyOrder.Filled.String()),
				sdk.NewAttribute(types.AttributeRemainingQuantity, buyOrder.Amount.Sub(buyOrder.Filled).String()),
				sdk.NewAttribute(types.AttributeOrderStatus, getOrderStatus(buyOrder)),
				sdk.NewAttribute("trades_executed", fmt.Sprintf("%d", tradesExecuted)),
				sdk.NewAttribute("total_volume", totalTradeVolume.String()),
				sdk.NewAttribute("total_volume_value", totalVolumeValue.String()),
			),
		)
	}

	return
}

func getOrderStatus(order *types.PerpetualOrder) string {
	if order.Filled.Equal(order.Amount) {
		return "filled"
	} else if order.Filled.IsPositive() {
		return "partially_filled"
	}
	return "open"
}
