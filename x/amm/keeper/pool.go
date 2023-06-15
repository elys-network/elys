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
	store.Set(types.PoolKey(
		pool.PoolId,
	), b)
	return nil
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(
	ctx sdk.Context,
	poolId uint64,

) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(
		poolId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(
	ctx sdk.Context,
	poolId uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(
		poolId,
	))
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
		return 0
	}
	return latestPool.PoolId + 1
}

// PoolExists checks if a pool with the given poolId exists in the list of pools
func (k Keeper) PoolExists(ctx sdk.Context, poolId uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	b := store.Get(types.PoolKey(poolId))
	return b != nil
}

// Get pool Ids that contains the denom in pool assets
func (k Keeper) GetAllPoolIdsWithDenom(ctx sdk.Context, denom string) (list []uint64) {
	pools := k.GetAllPool(ctx)

	for _, p := range pools {
		for _, asset := range p.PoolAssets {
			if denom == asset.Token.Denom {
				list = append(list, p.PoolId)
				break
			}
		}
	}

	return list
}
