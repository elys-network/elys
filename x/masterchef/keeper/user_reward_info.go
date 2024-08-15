package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetUserRewardInfo(ctx sdk.Context, userReward types.UserRewardInfo) {
	//store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&userReward)
	//key := types.GetUserRewardInfoKey(userReward.GetUserAccount(), userReward.GetPoolId(), userReward.GetRewardDenom())
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyUserRewardInfoKeyPrefix))
	key := types.LegacyUserRewardInfoKey(userReward.User, userReward.PoolId, userReward.RewardDenom)
	store.Set(key, b)
}

func (k Keeper) GetUserRewardInfo(ctx sdk.Context, user sdk.AccAddress, poolId uint64, rewardDenom string) (val types.UserRewardInfo, found bool) {
	//store := ctx.KVStore(k.storeKey)
	//key := types.GetUserRewardInfoKey(user, poolId, rewardDenom)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyUserRewardInfoKeyPrefix))
	key := types.LegacyUserRewardInfoKey(user.String(), poolId, rewardDenom)
	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

//func (k Keeper) RemoveUserRewardInfo(ctx sdk.Context, user sdk.AccAddress, poolId uint64, rewardDenom string) {
//	store := ctx.KVStore(k.storeKey)
//	key := types.GetUserRewardInfoKey(user, poolId, rewardDenom)
//	store.Delete(key)
//}

func (k Keeper) GetAllUserRewardInfos(ctx sdk.Context) (list []types.UserRewardInfo) {
	//store := ctx.KVStore(k.storeKey)
	//iterator := sdk.KVStorePrefixIterator(store, types.UserRewardInfoKeyPrefix)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyUserRewardInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UserRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// remove after migration
func (k Keeper) RemoveUserRewardInfo(ctx sdk.Context, user string, poolId uint64, rewardDenom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyUserRewardInfoKeyPrefix))
	store.Delete(types.LegacyUserRewardInfoKey(user, poolId, rewardDenom))
}

// remove after migration
//func (k Keeper) MigrateFromV2UserRewardInfos(ctx sdk.Context) {
//	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyUserRewardInfoKeyPrefix))
//	iterator := sdk.KVStorePrefixIterator(store, []byte{})
//
//	defer iterator.Close()
//
//	for ; iterator.Valid(); iterator.Next() {
//		var legacyUserRewardInfo types.UserRewardInfo
//		k.cdc.MustUnmarshal(iterator.Value(), &legacyUserRewardInfo)
//
//		k.SetUserRewardInfo(ctx, legacyUserRewardInfo)
//		k.DeleteLegacyUserRewardInfo(ctx, legacyUserRewardInfo.User, legacyUserRewardInfo.PoolId, legacyUserRewardInfo.RewardDenom)
//	}
//
//	return
//}
