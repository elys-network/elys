package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetAndUpdatePerpetualCounter(ctx sdk.Context, marketId uint64) uint64 {
	key := types.GetPerpetualCounterKey(marketId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		v := types.PerpetualCounter{
			MarketId: marketId,
			Counter:  2,
		}
		b = k.cdc.MustMarshal(&v)
		store.Set(key, b)
		return 1
	}

	var v types.PerpetualCounter
	k.cdc.MustUnmarshal(b, &v)
	result := v.Counter
	v.Counter = v.Counter + 1
	b = k.cdc.MustMarshal(&v)
	store.Set(key, b)
	return result
}

func (k Keeper) GetAllPerpetualCounters(ctx sdk.Context) []types.Perpetual {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualCounterPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.Perpetual

	for ; iterator.Valid(); iterator.Next() {
		var val types.Perpetual
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) setPerpetualCounter(ctx sdk.Context, p types.PerpetualCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualCounterKey(p.MarketId)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}
