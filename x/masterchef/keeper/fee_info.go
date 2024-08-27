// package keeper

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/elys-network/elys/x/masterchef/types"
// )

// func (k Keeper) SetFeeInfo(ctx sdk.Context, poolReward types.FeeInfo) {
// 	store := ctx.KVStore(k.storeKey)
// 	key := types.GetFeeInfoKey(poolReward.PoolId, poolReward.RewardDenom)
// 	b := k.cdc.MustMarshal(&poolReward)
// 	store.Set(key, b)
// }

// func (k Keeper) GetFeeInfo(ctx sdk.Context, poolId uint64, rewardDenom string) (val types.FeeInfo, found bool) {
// 	store := ctx.KVStore(k.storeKey)
// 	key := types.GetFeeInfoKey(poolId, rewardDenom)

// 	b := store.Get(key)
// 	if b == nil {
// 		return val, false
// 	}

// 	k.cdc.MustUnmarshal(b, &val)
// 	return val, true
// }

// func (k Keeper) RemoveFeeInfo(ctx sdk.Context, poolId uint64, rewardDenom string) {
// 	store := ctx.KVStore(k.storeKey)
// 	key := types.GetFeeInfoKey(poolId, rewardDenom)
// 	store.Delete(key)
// }

// func (k Keeper) GetAllFeeInfos(ctx sdk.Context) (list []types.FeeInfo) {
// 	store := ctx.KVStore(k.storeKey)
// 	iterator := sdk.KVStorePrefixIterator(store, types.FeeInfoKeyPrefix)

// 	defer iterator.Close()

// 	for ; iterator.Valid(); iterator.Next() {
// 		var val types.FeeInfo
// 		k.cdc.MustUnmarshal(iterator.Value(), &val)
// 		list = append(list, val)
// 	}

// 	return
// }
