package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) GetAndIncrementPerpetualCounter(ctx sdk.Context, marketId uint64) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketCounterKey(marketId)
	b := store.Get(key)
	var v types.PerpetualMarketCounter

	if b == nil {
		v = types.PerpetualMarketCounter{
			MarketId:          marketId,
			OrderCounter:      0,
			PerpetualCounter:  0,
			TotalOpenPosition: 0,
			TotalOpenOrders:   0,
		}
	}
	k.cdc.MustUnmarshal(b, &v)

	v.PerpetualCounter++
	v.TotalOpenPosition++

	store.Set(key, k.cdc.MustMarshal(&v))

	return v.PerpetualCounter
}

func (k Keeper) GetAndIncrementOrderCounter(ctx sdk.Context, marketId uint64) uint64 {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketCounterKey(marketId)
	b := store.Get(key)
	var v types.PerpetualMarketCounter

	if b == nil {
		v = types.PerpetualMarketCounter{
			MarketId:          marketId,
			OrderCounter:      0,
			PerpetualCounter:  0,
			TotalOpenPosition: 0,
			TotalOpenOrders:   0,
		}
	}
	k.cdc.MustUnmarshal(b, &v)

	v.OrderCounter++
	v.TotalOpenOrders++

	store.Set(key, k.cdc.MustMarshal(&v))

	return v.OrderCounter
}

func (k Keeper) DecrementTotalOrderCounter(ctx sdk.Context, marketId uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketCounterKey(marketId)
	b := store.Get(key)
	var v types.PerpetualMarketCounter

	if b == nil {
		return
	}
	k.cdc.MustUnmarshal(b, &v)

	v.TotalOpenOrders--
	store.Set(key, k.cdc.MustMarshal(&v))
}

func (k Keeper) DecrementTotalOpenPosition(ctx sdk.Context, marketId uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketCounterKey(marketId)
	b := store.Get(key)
	var v types.PerpetualMarketCounter

	if b == nil {
		return
	}
	k.cdc.MustUnmarshal(b, &v)

	v.TotalOpenPosition--
	store.Set(key, k.cdc.MustMarshal(&v))
}

func (k Keeper) SetPerpetualMarketCounter(ctx sdk.Context, v types.PerpetualMarketCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketCounterKey(v.MarketId)
	store.Set(key, k.cdc.MustMarshal(&v))
}

func (k Keeper) GetAllPerpetualMarketCounter(ctx sdk.Context) []types.PerpetualMarketCounter {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketCounterPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualMarketCounter

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualMarketCounter
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}
