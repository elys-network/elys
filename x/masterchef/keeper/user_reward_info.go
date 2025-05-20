package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/masterchef/types"
)

func (k Keeper) SetUserRewardInfo(ctx sdk.Context, userReward types.UserRewardInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&userReward)
	key := types.GetUserRewardInfoKey(userReward.GetUserAccount(), userReward.GetPoolId(), userReward.GetRewardDenom())
	store.Set(key, b)
}

func (k Keeper) GetUserRewardInfo(ctx sdk.Context, user sdk.AccAddress, poolId uint64, rewardDenom string) (val types.UserRewardInfo, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetUserRewardInfoKey(user, poolId, rewardDenom)
	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveUserRewardInfo(ctx sdk.Context, user sdk.AccAddress, poolId uint64, rewardDenom string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetUserRewardInfoKey(user, poolId, rewardDenom)
	store.Delete(key)
}

func (k Keeper) GetAllUserRewardInfos(ctx sdk.Context) (list []types.UserRewardInfo) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.UserRewardInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UserRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
