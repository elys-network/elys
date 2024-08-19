package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetUserRewardInfo(ctx sdk.Context, userReward types.UserRewardInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&userReward)
	key := types.GetUserRewardInfoKey(userReward.GetUserAccount(), userReward.GetPoolId(), userReward.GetRewardDenom())
	store.Set(key, b)
}

func (k Keeper) GetUserRewardInfo(ctx sdk.Context, user sdk.AccAddress, poolId uint64, rewardDenom string) (val types.UserRewardInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserRewardInfoKey(user, poolId, rewardDenom)
	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveUserRewardInfo(ctx sdk.Context, user sdk.AccAddress, poolId uint64, rewardDenom string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUserRewardInfoKey(user, poolId, rewardDenom)
	store.Delete(key)
}

func (k Keeper) GetAllUserRewardInfos(ctx sdk.Context) (list []types.UserRewardInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UserRewardInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UserRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// remove after migration
func (k Keeper) MigrateFromV3UserRewardInfos(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.LegacyUserRewardInfoKeyPrefix))

	defer iterator.Close()

	keysToDelete := [][]byte{}
	for ; iterator.Valid(); iterator.Next() {
		var userRewardInfo types.UserRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &userRewardInfo)

		if !userRewardInfo.RewardPending.IsZero() || !userRewardInfo.RewardDebt.IsZero() {
			k.SetUserRewardInfo(ctx, userRewardInfo)
		}

		key := types.GetLegacyUserRewardInfoKey(userRewardInfo.User, userRewardInfo.PoolId, userRewardInfo.RewardDenom)
		keysToDelete = append(keysToDelete, key)

		if len(keysToDelete) == 10_000 {
			k.deleteLegacyUserRewardInfos(ctx, keysToDelete)
			keysToDelete = [][]byte{}
		}
	}

	k.deleteLegacyUserRewardInfos(ctx, keysToDelete)
	return
}

func (k Keeper) deleteLegacyUserRewardInfos(ctx sdk.Context, keysToDelete [][]byte) {
	store := ctx.KVStore(k.storeKey)
	for _, key := range keysToDelete {
		store.Delete(key)
	}

	return
}
