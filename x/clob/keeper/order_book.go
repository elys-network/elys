package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

// GetBuyOrderBookDepth returns the number of buy orders for a market
func (k Keeper) GetBuyOrderBookDepth(ctx sdk.Context, marketId uint64) uint64 {
	iterator := k.GetBuyOrderIterator(ctx, marketId)
	defer iterator.Close()

	var count uint64
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return count
}

// GetSellOrderBookDepth returns the number of sell orders for a market
func (k Keeper) GetSellOrderBookDepth(ctx sdk.Context, marketId uint64) uint64 {
	iterator := k.GetSellOrderIterator(ctx, marketId)
	defer iterator.Close()

	var count uint64
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return count
}

// GetOrderBookSnapshot returns a snapshot of the order book
func (k Keeper) GetOrderBookSnapshot(ctx sdk.Context, marketId uint64, maxDepth uint32) types.OrderBookSnapshot {
	buyOrders := k.GetBuyOrdersUpToDepth(ctx, marketId, maxDepth)
	sellOrders := k.GetSellOrdersUpToDepth(ctx, marketId, maxDepth)

	return types.OrderBookSnapshot{
		MarketId:   marketId,
		BuyOrders:  buyOrders,
		SellOrders: sellOrders,
		Timestamp:  uint64(ctx.BlockTime().Unix()),
	}
}

// GetBuyOrdersUpToDepth returns buy orders up to specified depth (highest price first)
func (k Keeper) GetBuyOrdersUpToDepth(ctx sdk.Context, marketId uint64, maxDepth uint32) []types.OrderBookEntry {
	// Get all buy orders first since we need them in reverse order
	iterator := k.GetBuyOrderIterator(ctx, marketId)
	defer iterator.Close()

	var allOrders []types.PerpetualOrder
	for ; iterator.Valid(); iterator.Next() {
		var order types.PerpetualOrder
		if err := k.cdc.Unmarshal(iterator.Value(), &order); err != nil {
			ctx.Logger().Error("failed to unmarshal buy order in GetBuyOrdersUpToDepth", "error", err)
			continue
		}
		allOrders = append(allOrders, order)
	}

	// Process in reverse order (highest price first) up to maxDepth
	var orders []types.OrderBookEntry
	count := uint32(0)
	for i := len(allOrders) - 1; i >= 0 && count < maxDepth; i-- {
		order := allOrders[i]
		orders = append(orders, types.OrderBookEntry{
			Price:    order.GetPrice(),
			Quantity: order.Amount.Sub(order.Filled),
		})
		count++
	}

	return orders
}

// GetSellOrdersUpToDepth returns sell orders up to specified depth (lowest price first)
func (k Keeper) GetSellOrdersUpToDepth(ctx sdk.Context, marketId uint64, maxDepth uint32) []types.OrderBookEntry {
	iterator := k.GetSellOrderIterator(ctx, marketId)
	defer iterator.Close()

	var orders []types.OrderBookEntry
	var count uint32

	for ; iterator.Valid() && count < maxDepth; iterator.Next() {
		var order types.PerpetualOrder
		if err := k.cdc.Unmarshal(iterator.Value(), &order); err != nil {
			ctx.Logger().Error("failed to unmarshal sell order in GetSellOrdersUpToDepth", "error", err)
			continue
		}

		orders = append(orders, types.OrderBookEntry{
			Price:    order.GetPrice(),
			Quantity: order.Amount.Sub(order.Filled),
		})
		count++
	}

	return orders
}
