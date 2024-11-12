package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

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

func (k Keeper) DeleteLegacyPoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.LegacyPoolRewardInfoKeyPrefix))
	store.Delete(types.LegacyPoolRewardInfoKey(poolId, rewardDenom))
}

func (k Keeper) GetAllLegacyPoolRewardInfos(ctx sdk.Context) (list []types.PoolRewardInfo) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.LegacyPoolRewardInfoKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
