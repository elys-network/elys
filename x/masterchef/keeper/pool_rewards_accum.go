package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) GetPoolRewardsAccum(ctx sdk.Context, poolId, timestamp uint64) (types.PoolRewardsAccum, error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetPoolRewardsAccumKey(poolId, timestamp))
	if b == nil {
		return types.PoolRewardsAccum{}, types.ErrPoolRewardsAccumNotFound
	}

	accum := types.PoolRewardsAccum{}
	k.cdc.MustUnmarshal(b, &accum)
	return accum, nil
}

func (k Keeper) SetPoolRewardsAccum(ctx sdk.Context, accum types.PoolRewardsAccum) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&accum)
	store.Set(types.GetPoolRewardsAccumKey(accum.PoolId, accum.Timestamp), bz)
}

func (k Keeper) DeletePoolRewardsAccum(ctx sdk.Context, accum types.PoolRewardsAccum) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolRewardsAccumKey(accum.PoolId, accum.Timestamp))
}

func (k Keeper) GetAllPoolRewardsAccum(ctx sdk.Context) (list []types.PoolRewardsAccum) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolRewardsAccumKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardsAccum
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) IterateAllPoolRewardsAccum(ctx sdk.Context, handler func(accum types.PoolRewardsAccum) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte(types.PoolRewardsAccumKeyPrefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		accum := types.PoolRewardsAccum{}
		k.cdc.MustUnmarshal(iter.Value(), &accum)
		if handler(accum) {
			break
		}
	}
}

func (k Keeper) IteratePoolRewardsAccum(ctx sdk.Context, poolId uint64, handler func(accum types.PoolRewardsAccum) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		accum := types.PoolRewardsAccum{}
		k.cdc.MustUnmarshal(iter.Value(), &accum)
		if handler(accum) {
			break
		}
	}
}

func (k Keeper) FirstPoolRewardsAccum(ctx sdk.Context, poolId uint64) types.PoolRewardsAccum {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		accum := types.PoolRewardsAccum{}
		k.cdc.MustUnmarshal(iter.Value(), &accum)
		return accum
	}
	return types.PoolRewardsAccum{}
}

func (k Keeper) LastPoolRewardsAccum(ctx sdk.Context, poolId uint64) types.PoolRewardsAccum {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		accum := types.PoolRewardsAccum{}
		k.cdc.MustUnmarshal(iter.Value(), &accum)
		return accum
	}
	return types.PoolRewardsAccum{
		PoolId:      poolId,
		BlockHeight: 0,
		Timestamp:   0,
		DexReward:   math.LegacyZeroDec(),
		GasReward:   math.LegacyZeroDec(),
		EdenReward:  math.LegacyZeroDec(),
	}
}

func (k Keeper) AddPoolRewardsAccum(ctx sdk.Context, poolId, timestamp uint64, height int64, dexReward, gasReward, edenReward math.LegacyDec) {
	lastAccum := k.LastPoolRewardsAccum(ctx, poolId)
	lastAccum.Timestamp = timestamp
	lastAccum.BlockHeight = height
	if lastAccum.DexReward.IsNil() {
		lastAccum.DexReward = math.LegacyZeroDec()
	}
	if lastAccum.GasReward.IsNil() {
		lastAccum.GasReward = math.LegacyZeroDec()
	}
	if lastAccum.EdenReward.IsNil() {
		lastAccum.EdenReward = math.LegacyZeroDec()
	}
	lastAccum.DexReward = lastAccum.DexReward.Add(dexReward)
	lastAccum.GasReward = lastAccum.GasReward.Add(gasReward)
	lastAccum.EdenReward = lastAccum.EdenReward.Add(edenReward)
	k.SetPoolRewardsAccum(ctx, lastAccum)
}
