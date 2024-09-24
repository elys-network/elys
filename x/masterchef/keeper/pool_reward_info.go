package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetPoolRewardInfo(ctx sdk.Context, poolReward types.PoolRewardInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolRewardInfoKey(poolReward.PoolId, poolReward.RewardDenom)
	b := k.cdc.MustMarshal(&poolReward)
	store.Set(key, b)
}

func (k Keeper) GetPoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) (val types.PoolRewardInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolRewardInfoKey(poolId, rewardDenom)

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolRewardInfoKey(poolId, rewardDenom)
	store.Delete(key)
}

func (k Keeper) GetAllPoolRewardInfos(ctx sdk.Context) (list []types.PoolRewardInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolRewardInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) DeleteLegacyPoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyPoolRewardInfoKeyPrefix))
	store.Delete(types.LegacyPoolRewardInfoKey(poolId, rewardDenom))
}

func (k Keeper) GetAllLegacyPoolRewardInfos(ctx sdk.Context) (list []types.PoolRewardInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyPoolRewardInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
