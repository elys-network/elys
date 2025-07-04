package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (k Keeper) GetAllPositionCounters(ctx sdk.Context) []types.PositionCounter {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PositionCounterPrefix)
	defer iterator.Close()

	var list []types.PositionCounter

	for ; iterator.Valid(); iterator.Next() {
		var positionCounter types.PositionCounter
		k.cdc.MustUnmarshal(iterator.Value(), &positionCounter)
		list = append(list, positionCounter)
	}

	return list
}

func (k Keeper) GetPositionCounter(ctx sdk.Context, poolId uint64) types.PositionCounter {
	var positionCounter types.PositionCounter
	key := types.GetPositionCounterKey(poolId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(key)
	if bz == nil {
		return types.PositionCounter{
			PoolId:    poolId,
			Counter:   0,
			TotalOpen: 0,
		}
	}
	k.cdc.MustUnmarshal(bz, &positionCounter)
	return positionCounter
}

func (k Keeper) SetPositionCounter(ctx sdk.Context, positionCounter types.PositionCounter) {
	key := types.GetPositionCounterKey(positionCounter.PoolId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&positionCounter)
	store.Set(key, bz)
}
