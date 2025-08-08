package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) GetPerpetualOrder(ctx sdk.Context, orderId types.OrderId) (types.Order, bool) {
	key := types.GetOrderKey(orderId.MarketId, orderId.OrderType, orderId.PriceTick, orderId.Counter)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.Order{}, false
	}

	var val types.Order
	if err := k.cdc.Unmarshal(b, &val); err != nil {
		ctx.Logger().Error("failed to unmarshal perpetual order", "key", orderId, "error", err)
		return types.Order{}, false
	}
	return val, true
}

func (k Keeper) GetAllOrders(ctx sdk.Context) []types.Order {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOrderPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.Order

	for ; iterator.Valid(); iterator.Next() {
		var val types.Order
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetOrder(ctx sdk.Context, v types.Order) error {
	if v.GetOrderType() == types.OrderType_ORDER_TYPE_LIMIT_SELL || v.GetOrderType() == types.OrderType_ORDER_TYPE_LIMIT_BUY || v.GetPriceTick() <= 0 {
		store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		key := types.GetOrderKey(v.GetMarketId(), v.GetOrderType(), v.GetPriceTick(), v.GetCounter())
		b := k.cdc.MustMarshal(&v)
		store.Set(key, b)

		// Sync with memory orderbook
		// Check if this is a new order or an update
		if v.Filled.IsZero() {
			// New order - add to memory
			k.memoryOrderBook.AddOrder(v.GetMarketId(), &v)
		} else {
			// Update existing order in memory
			k.memoryOrderBook.UpdateOrder(v.GetMarketId(), &v)
		}

		return nil
	} else {
		return types.ErrInvalidOrderType
	}
}

func (k Keeper) DeleteOrder(ctx sdk.Context, perpetualOrderOwner types.PerpetualOrderOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetOrderKey(perpetualOrderOwner.OrderId.MarketId, perpetualOrderOwner.OrderId.OrderType, perpetualOrderOwner.OrderId.PriceTick, perpetualOrderOwner.OrderId.Counter)
	store.Delete(key)

	k.DeleteOrderOwner(ctx, perpetualOrderOwner)

	// Remove from memory orderbook
	k.memoryOrderBook.RemoveOrder(
		perpetualOrderOwner.OrderId.MarketId,
		perpetualOrderOwner.OrderId.Counter,
		types.IsBuy(perpetualOrderOwner.OrderId.OrderType),
	)
}

func (k Keeper) GetBuyOrderIterator(ctx sdk.Context, marketId uint64) storetypes.Iterator {
	key := types.GetOrderBookIteratorKey(marketId, true)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	return storetypes.KVStoreReversePrefixIterator(store, []byte{})

}

func (k Keeper) GetSellOrderIterator(ctx sdk.Context, marketId uint64) storetypes.Iterator {
	key := types.GetOrderBookIteratorKey(marketId, false)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	return storetypes.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) RequiredBalanceForOrder(ctx sdk.Context, order types.Order) (sdk.Coin, error) {
	market, err := k.GetPerpetualMarket(ctx, order.GetMarketId())
	if err != nil {
		return sdk.Coin{}, err
	}

	maxFees := k.CalculateMaxFees(market, order)
	initialMarginRequired := order.UnfilledValue().Mul(market.InitialMarginRatio).RoundInt()

	amount := maxFees.Add(initialMarginRequired)

	return sdk.NewCoin(market.QuoteDenom, amount), nil
}
