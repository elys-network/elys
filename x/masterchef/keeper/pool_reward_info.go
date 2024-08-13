package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetPoolRewardInfo(ctx sdk.Context, poolReward types.PoolRewardInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolRewardInfoKeyPrefix))
	b := k.cdc.MustMarshal(&poolReward)
	store.Set(types.PoolRewardInfoKey(poolReward.PoolId, poolReward.RewardDenom), b)
}

func (k Keeper) GetPoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) (val types.PoolRewardInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolRewardInfoKeyPrefix))

	b := store.Get(types.PoolRewardInfoKey(poolId, rewardDenom))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePoolRewardInfo(ctx sdk.Context, poolId uint64, rewardDenom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolRewardInfoKeyPrefix))
	store.Delete(types.PoolRewardInfoKey(poolId, rewardDenom))
}

func (k Keeper) GetAllPoolRewardInfos(ctx sdk.Context) (list []types.PoolRewardInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolRewardInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
