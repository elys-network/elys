package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k Keeper) SetPoolInfo(ctx sdk.Context, pool types.PoolInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolInfoKey(pool.PoolId)
	b := k.cdc.MustMarshal(&pool)
	store.Set(key, b)
}

func (k Keeper) GetPoolInfo(ctx sdk.Context, poolId uint64) (val types.PoolInfo, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolInfoKey(poolId)

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePoolInfo(ctx sdk.Context, poolId uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolInfoKey(poolId)
	store.Delete(key)
}

func (k Keeper) GetAllPoolInfos(ctx sdk.Context) (list []types.PoolInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetPoolRewardInfo(ctx sdk.Context, poolReward types.PoolRewardInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolRewardInfoKey(poolReward.PoolId, poolReward.RewardDenom)
	b := k.cdc.MustMarshal(&poolReward)
	store.Set(key, b)
}

func (k Keeper) GetPoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) (val types.PoolRewardInfo, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolRewardInfoKey(poolId, rewardDenom)

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolRewardInfoKey(poolId, rewardDenom)
	store.Delete(key)
}

func (k Keeper) GetAllPoolRewardInfos(ctx sdk.Context) (list []types.PoolRewardInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolRewardInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
