package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetPoolRewardsAccum(ctx sdk.Context, poolId, timestamp uint64) (types.PoolRewardsAccum, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := store.Get(types.GetPoolRewardsAccumKey(poolId, timestamp))
	if b == nil {
		return types.PoolRewardsAccum{}, types.ErrPoolRewardsAccumNotFound
	}

	accum := types.PoolRewardsAccum{}
	k.cdc.MustUnmarshal(b, &accum)
	return accum, nil
}

func (k Keeper) SetPoolRewardsAccum(ctx sdk.Context, accum types.PoolRewardsAccum) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&accum)
	store.Set(types.GetPoolRewardsAccumKey(accum.PoolId, accum.Timestamp), bz)
}

func (k Keeper) DeletePoolRewardsAccum(ctx sdk.Context, accum types.PoolRewardsAccum) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetPoolRewardsAccumKey(accum.PoolId, accum.Timestamp))
}

func (k Keeper) GetAllPoolRewardsAccum(ctx sdk.Context) (list []types.PoolRewardsAccum) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolRewardsAccumKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolRewardsAccum
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) IterateAllPoolRewardsAccum(ctx sdk.Context, handler func(accum types.PoolRewardsAccum) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.PoolRewardsAccumKeyPrefix)
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStorePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		accum := types.PoolRewardsAccum{}
		k.cdc.MustUnmarshal(iter.Value(), &accum)
		return accum
	}
	return types.PoolRewardsAccum{}
}

func (k Keeper) LastPoolRewardsAccum(ctx sdk.Context, poolId uint64) types.PoolRewardsAccum {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStoreReversePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
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

// Returns eden rewards using forward calc for 24 hours
func (k Keeper) ForwardEdenCalc(ctx sdk.Context, poolId uint64) osmomath.BigDec {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iter := storetypes.KVStoreReversePrefixIterator(store, types.GetPoolRewardsAccumPrefix(poolId))
	defer iter.Close()

	var lastTwo []types.PoolRewardsAccum
	for ; iter.Valid() && len(lastTwo) < 2; iter.Next() {
		accum := types.PoolRewardsAccum{}
		k.cdc.MustUnmarshal(iter.Value(), &accum)
		lastTwo = append([]types.PoolRewardsAccum{accum}, lastTwo...)
	}

	if len(lastTwo) == 2 {
		diff := lastTwo[1].GetBigDecEdenReward().Sub(lastTwo[0].GetBigDecEdenReward())
		// Here we are assuming average block time of 4s
		// 1 DAY = 86400
		// Note: This calculation maybe used in FE, the idea is to
		// give estimated numbers of rewards that a user will get
		return diff.MulInt64(21600)
	}

	// Return zero if there are not enough entries
	return osmomath.ZeroBigDec()
}

func (k Keeper) AddPoolRewardsAccum(ctx sdk.Context, poolId, timestamp uint64, height int64, dexReward, gasReward, edenReward osmomath.BigDec) {
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
	lastAccum.DexReward = lastAccum.GetBigDecDexReward().Add(dexReward).Dec()
	lastAccum.GasReward = lastAccum.GetBigDecGasReward().Add(gasReward).Dec()
	lastAccum.EdenReward = lastAccum.GetBigDecEdenReward().Add(edenReward).Dec()
	k.SetPoolRewardsAccum(ctx, lastAccum)
}
