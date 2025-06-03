package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetualOrder(ctx sdk.Context, marketId uint64, orderType types.OrderType, price math.LegacyDec, blockHeight uint64) (types.PerpetualOrder, bool) {
	key := types.GetPerpetualOrderKey(marketId, orderType, price, blockHeight)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.PerpetualOrder{}, false
	}

	var val types.PerpetualOrder
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllPerpetualOrders(ctx sdk.Context) []types.PerpetualOrder {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOrderPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualOrder

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetualOrder(ctx sdk.Context, v types.PerpetualOrder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOrderKey(v.MarketId, v.OrderType, v.Price, v.Counter)
	b := k.cdc.MustMarshal(&v)
	store.Set(key, b)
}

func (k Keeper) DeleteOrder(ctx sdk.Context, perpetualOrderOwner types.PerpetualOrderOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOrderKey(perpetualOrderOwner.OrderKey.MarketId, perpetualOrderOwner.OrderKey.OrderType, perpetualOrderOwner.OrderKey.Price, perpetualOrderOwner.OrderKey.Counter)
	store.Delete(key)

	k.DeleteOrderOwner(ctx, perpetualOrderOwner)
}

func (k Keeper) GetOrderBookWithSide(ctx sdk.Context, marketId uint64, long bool, count uint64) []types.PerpetualOrder {
	key := types.GetPerpetualOrderBookIteratorKey(marketId, long)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	var list []types.PerpetualOrder

	for ; iterator.Valid() && count > 0; iterator.Next() {
		var val types.PerpetualOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
		count--
	}
	return list
}

func (k Keeper) GetBuyOrderIterator(ctx sdk.Context, marketId uint64) storetypes.Iterator {
	key := types.GetPerpetualOrderBookIteratorKey(marketId, true)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	return storetypes.KVStoreReversePrefixIterator(store, []byte{})

}

func (k Keeper) GetSellOrderIterator(ctx sdk.Context, marketId uint64) storetypes.Iterator {
	key := types.GetPerpetualOrderBookIteratorKey(marketId, false)
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), key)
	return storetypes.KVStorePrefixIterator(store, []byte{})
}
