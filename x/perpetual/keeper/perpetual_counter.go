package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) GetPerpetualCounter(ctx sdk.Context, poolId uint64) (v types.PerpetualCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GePerpetualCounterKey(poolId)
	bz := store.Get(key)
	if bz == nil {
		v = types.PerpetualCounter{
			AmmPoolId: poolId,
			Counter:   0,
			TotalOpen: 0,
		}
		k.SetPerpetualCounter(ctx, v)
		return
	}
	k.cdc.MustUnmarshal(bz, &v)
	return
}

func (k Keeper) GetAllPerpetualCounter(ctx sdk.Context) (list []types.PerpetualCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PerpetualCounterPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var v types.PerpetualCounter
		bytesValue := iterator.Value()
		k.cdc.MustUnmarshal(bytesValue, &v)
		list = append(list, v)
	}
	return list
}

func (k Keeper) SetPerpetualCounter(ctx sdk.Context, v types.PerpetualCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GePerpetualCounterKey(v.AmmPoolId)
	store.Set(key, k.cdc.MustMarshal(&v))
}
