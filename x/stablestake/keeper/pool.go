package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, types.PoolPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pool := types.Pool{}
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		pools = append(pools, pool)
	}

	return
}

// GetPool get pool as types.Pool
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (pool types.Pool, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetPoolKey(poolId))
	if b == nil {
		return types.Pool{}, false
	}

	k.cdc.MustUnmarshal(b, &pool)
	return pool, true
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.GetPoolKey(pool.Id), b)
}

func (k Keeper) DeletePool(ctx sdk.Context, poolId uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetPoolKey(poolId))
}

func (k Keeper) GetPoolByDenom(ctx sdk.Context, denom string) (types.Pool, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, types.PoolPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pool := types.Pool{}
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if pool.DepositDenom == denom {
			return pool, true
		}
	}
	return types.Pool{}, false
}

func (k Keeper) CalculateRedemptionRateForPool(ctx sdk.Context, pool types.Pool) osmomath.BigDec {
	totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForPool(pool.Id))

	if totalShares.Amount.IsZero() {
		return osmomath.ZeroBigDec()
	}

	return pool.GetBigDecNetAmount().Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))
}

func (k Keeper) GetLatestPool(ctx sdk.Context) (val types.Pool, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStoreReversePrefixIterator(store, types.PoolPrefixKey)
	defer iterator.Close()

	if !iterator.Valid() {
		return val, false
	}

	k.cdc.MustUnmarshal(iterator.Value(), &val)
	return val, true
}

// GetNextPoolId returns the next pool id.
func (k Keeper) GetNextPoolId(ctx sdk.Context) uint64 {
	latestPool, found := k.GetLatestPool(ctx)
	if !found {
		return types.UsdcPoolId
	}
	return latestPool.Id + 1
}

// IterateLiquidityPools iterates over all LiquidityPools and performs a
// callback.
func (k Keeper) IterateLiquidityPools(ctx sdk.Context, handlerFn func(pool types.Pool) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PoolPrefixKey)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if handlerFn(pool) {
			break
		}
	}
}

func (k Keeper) HasPoolByDenom(ctx sdk.Context, depositDenom string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, types.PoolPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pool := types.Pool{}
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if pool.DepositDenom == depositDenom {
			return true
		}
	}

	return false
}
