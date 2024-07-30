package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, index uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(index))
}

// GetAllPool returns all pool
func (k Keeper) GetAllPools(ctx sdk.Context) (list []types.Pool) {
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

// SetPool set a specific pool in the store from its index
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKey(pool.AmmPoolId), b)
}

func (k Keeper) DeletePool(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(poolId))
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
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

func (k Keeper) GetEnabledPools(ctx sdk.Context) []uint64 {
	poolIds := make([]uint64, 0)
	pools := k.GetAllPools(ctx)
	for _, p := range pools {
		if p.Enabled {
			poolIds = append(poolIds, p.AmmPoolId)
		}
	}

	return poolIds
}

func (k Keeper) SetEnabledPools(ctx sdk.Context, pools []uint64) {
	for _, poolId := range pools {
		pool, found := k.GetPool(ctx, poolId)
		if !found {
			pool = types.NewPool(poolId)
			k.SetPool(ctx, pool)
		}
		pool.Enabled = true

		k.SetPool(ctx, pool)
	}
}

func (k Keeper) IsPoolEnabled(ctx sdk.Context, poolId uint64) bool {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		pool = types.NewPool(poolId)
		k.SetPool(ctx, pool)
	}

	return pool.Enabled
}

func (k Keeper) IsPoolClosed(ctx sdk.Context, poolId uint64) bool {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		pool = types.NewPool(poolId)
		k.SetPool(ctx, pool)
	}

	return pool.Closed
}
