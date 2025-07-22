package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// BatchDeleteOrders deletes multiple orders efficiently in a single operation
func (k Keeper) BatchDeleteOrders(ctx sdk.Context, keys []types.PerpetualOrderOwner) {
	if len(keys) == 0 {
		return
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	// Delete all orders and their owners in batch
	for _, key := range keys {
		// Delete the order
		orderKey := types.GetPerpetualOrderKey(key.OrderKey.MarketId, key.OrderKey.OrderType, key.OrderKey.Price, key.OrderKey.Counter)
		store.Delete(orderKey)

		// Delete the order owner
		k.DeleteOrderOwner(ctx, key)
	}

	// Clear any cached market data for affected markets
	marketSet := make(map[uint64]bool)
	for _, key := range keys {
		marketSet[key.OrderKey.MarketId] = true
	}

	// Invalidate cache for affected markets
	if k.marketCache != nil {
		k.marketCache.mu.Lock()
		for marketId := range marketSet {
			delete(k.marketCache.bestBids, marketId)
			delete(k.marketCache.bestAsks, marketId)
		}
		k.marketCache.mu.Unlock()
	}
}

// BatchUpdatePerpetuals updates multiple perpetuals efficiently
func (k Keeper) BatchUpdatePerpetuals(ctx sdk.Context, perpetuals []types.Perpetual) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	for _, perpetual := range perpetuals {
		key := types.GetPerpetualKey(perpetual.MarketId, perpetual.Id)
		b, err := k.cdc.Marshal(&perpetual)
		if err != nil {
			return err
		}
		store.Set(key, b)

		// Update owner mapping
		ownerKey := types.GetPerpetualOwnerKey(perpetual.GetOwnerAccAddress(), perpetual.SubAccountId, perpetual.MarketId, perpetual.Id)
		owner := types.PerpetualOwner{
			Owner:        perpetual.Owner,
			SubAccountId: perpetual.SubAccountId,
			MarketId:     perpetual.MarketId,
			PerpetualId:  perpetual.Id,
		}
		ownerB, err := k.cdc.Marshal(&owner)
		if err != nil {
			return err
		}
		store.Set(ownerKey, ownerB)
	}

	return nil
}

// BatchGetPerpetuals retrieves multiple perpetuals efficiently
func (k Keeper) BatchGetPerpetuals(ctx sdk.Context, requests []struct {
	MarketId    uint64
	PerpetualId uint64
}) ([]types.Perpetual, error) {
	perpetuals := make([]types.Perpetual, 0, len(requests))

	for _, req := range requests {
		perpetual, err := k.GetPerpetual(ctx, req.MarketId, req.PerpetualId)
		if err != nil {
			ctx.Logger().Error("failed to get perpetual", "marketId", req.MarketId, "perpetualId", req.PerpetualId, "error", err)
			continue
		}
		perpetuals = append(perpetuals, perpetual)
	}

	return perpetuals, nil
}

// OptimizedOrderMatching performs order matching with reduced database operations
func (k Keeper) OptimizedOrderMatching(ctx sdk.Context, marketId uint64, maxMatches int) error {
	// Get best bid and ask prices from cache
	bestBid := k.GetCachedBestBid(ctx, marketId)
	bestAsk := k.GetCachedBestAsk(ctx, marketId)

	// Quick check if matching is possible
	if bestBid.IsZero() || bestAsk.IsZero() || bestBid.LT(bestAsk) {
		return nil // No matching possible
	}

	// Collect orders to match
	var buyOrders, sellOrders []types.PerpetualOrder
	var buyKeys, sellKeys []types.PerpetualOrderOwner
	matchCount := 0

	// Get buy orders
	buyIterator := k.GetBuyOrderIterator(ctx, marketId)
	defer buyIterator.Close()

	for ; buyIterator.Valid() && matchCount < maxMatches; buyIterator.Next() {
		var buyOrder types.PerpetualOrder
		if err := k.cdc.Unmarshal(buyIterator.Value(), &buyOrder); err != nil {
			ctx.Logger().Error("failed to unmarshal buy order", "error", err)
			continue
		}

		if buyOrder.Price.GTE(bestAsk) {
			buyOrders = append(buyOrders, buyOrder)
			buyKeys = append(buyKeys, types.PerpetualOrderOwner{
				Owner:        buyOrder.Owner,
				SubAccountId: buyOrder.SubAccountId,
				OrderKey:     types.NewOrderKey(buyOrder.MarketId, buyOrder.OrderType, buyOrder.Price, buyOrder.Counter),
			})
			matchCount++
		} else {
			break // No more matches possible
		}
	}

	// Get sell orders
	sellIterator := k.GetSellOrderIterator(ctx, marketId)
	defer sellIterator.Close()

	matchCount = 0
	for ; sellIterator.Valid() && matchCount < maxMatches; sellIterator.Next() {
		var sellOrder types.PerpetualOrder
		if err := k.cdc.Unmarshal(sellIterator.Value(), &sellOrder); err != nil {
			ctx.Logger().Error("failed to unmarshal sell order", "error", err)
			continue
		}

		if sellOrder.Price.LTE(bestBid) {
			sellOrders = append(sellOrders, sellOrder)
			sellKeys = append(sellKeys, types.PerpetualOrderOwner{
				Owner:        sellOrder.Owner,
				SubAccountId: sellOrder.SubAccountId,
				OrderKey:     types.NewOrderKey(sellOrder.MarketId, sellOrder.OrderType, sellOrder.Price, sellOrder.Counter),
			})
			matchCount++
		} else {
			break // No more matches possible
		}
	}

	// Match orders
	var ordersToDelete []types.PerpetualOrderOwner

	for i := 0; i < len(buyOrders) && i < len(sellOrders); i++ {
		buyOrder := &buyOrders[i]
		sellOrder := &sellOrders[i]

		// Determine trade price
		tradePrice := sellOrder.Price
		if buyOrder.Counter < sellOrder.Counter {
			tradePrice = buyOrder.Price
		}

		// Calculate trade quantity
		buyRemaining := buyOrder.Amount.Sub(buyOrder.Filled)
		sellRemaining := sellOrder.Amount.Sub(sellOrder.Filled)
		tradeQuantity := buyRemaining
		if sellRemaining.LT(buyRemaining) {
			tradeQuantity = sellRemaining
		}

		// Update filled amounts
		buyOrder.Filled = buyOrder.Filled.Add(tradeQuantity)
		sellOrder.Filled = sellOrder.Filled.Add(tradeQuantity)

		// Execute trade
		buyerSubAccount, err := k.GetSubAccount(ctx, buyOrder.GetOwnerAccAddress(), buyOrder.SubAccountId)
		if err != nil {
			continue
		}

		sellerSubAccount, err := k.GetSubAccount(ctx, sellOrder.GetOwnerAccAddress(), sellOrder.SubAccountId)
		if err != nil {
			continue
		}

		err = k.Exchange(ctx, types.Trade{
			BuyerSubAccount:     buyerSubAccount,
			SellerSubAccount:    sellerSubAccount,
			MarketId:            marketId,
			Price:               tradePrice,
			Quantity:            tradeQuantity,
			IsBuyerLiquidation:  false,
			IsSellerLiquidation: false,
			IsBuyerTaker:        buyOrder.Counter > sellOrder.Counter,
		})

		if err != nil {
			ctx.Logger().Error("failed to execute trade", "error", err)
			continue
		}

		// Mark fully filled orders for deletion
		if buyOrder.Filled.Equal(buyOrder.Amount) {
			ordersToDelete = append(ordersToDelete, buyKeys[i])
		} else {
			k.SetPerpetualOrder(ctx, *buyOrder)
		}

		if sellOrder.Filled.Equal(sellOrder.Amount) {
			ordersToDelete = append(ordersToDelete, sellKeys[i])
		} else {
			k.SetPerpetualOrder(ctx, *sellOrder)
		}
	}

	// Batch delete fully filled orders
	k.BatchDeleteOrders(ctx, ordersToDelete)

	return nil
}
