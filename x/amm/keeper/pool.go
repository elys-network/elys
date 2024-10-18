package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SetPool set a specific pool in the store from its index
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKey(pool.PoolId), b)
	return nil
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(poolId))
}

// GetAllPool returns all pool
func (k Keeper) GetAllPool(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllLegacyPool returns all legacy pool
func (k Keeper) GetAllLegacyPool(ctx sdk.Context) (list []types.LegacyPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetLatestPool retrieves the latest pool item from the list of pools
func (k Keeper) GetLatestPool(ctx sdk.Context) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStoreReversePrefixIterator(store, []byte{})
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
		return 1
	}
	return latestPool.PoolId + 1
}

// PoolExists checks if a pool with the given poolId exists in the list of pools
func (k Keeper) PoolExists(ctx sdk.Context, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	b := store.Get(types.PoolKey(poolId))
	return b != nil
}

// GetBestPoolWithDenoms returns the first highest TVL pool id that contains all specified denominations
func (k Keeper) GetBestPoolWithDenoms(ctx sdk.Context, denoms []string, usesOracle bool) (pool types.Pool, found bool) {
	// Get all pools
	pools := k.GetAllPool(ctx)

	maxTvl := sdk.NewDec(-1)
	bestPool := types.Pool{}
	for _, p := range pools {
		// If usesOracle is false, function filters in all pools.
		if usesOracle && !p.PoolParams.UseOracle {
			continue
		}
		// If the number of assets in the pool is less than the number of denoms, skip
		if len(p.PoolAssets) < len(denoms) {
			continue
		}

		allDenomsFound := true

		// Check that all denoms are in the pool's assets
		for _, denom := range denoms {
			denomFound := false
			for _, asset := range p.PoolAssets {
				if denom == asset.Token.Denom {
					denomFound = true
					break
				}
			}

			// If any denom is not found, mark allDenomsFound as false and break
			if !denomFound {
				allDenomsFound = false
				break
			}
		}

		poolTvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			poolTvl = sdk.ZeroDec()
		}

		// If all denoms are found in this pool, return the pool id
		if allDenomsFound && maxTvl.LT(poolTvl) {
			maxTvl = poolTvl
			bestPool = p
		}
	}

	return bestPool, !maxTvl.IsNegative()
}

// IterateLiquidty iterates over all LiquidityPools and performs a
// callback.
func (k Keeper) IterateLiquidityPools(ctx sdk.Context, handlerFn func(pool types.Pool) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if handlerFn(pool) {
			break
		}
	}

	return
}

// GetPoolSnapshotOrSet returns a pool snapshot or set the snapshot
func (k Keeper) GetPoolSnapshotOrSet(ctx sdk.Context, pool types.Pool) (val types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.transientStoreKey), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(pool.PoolId))
	if b == nil {
		b := k.cdc.MustMarshal(&pool)
		store.Set(types.PoolKey(pool.PoolId), b)
		return pool
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}
