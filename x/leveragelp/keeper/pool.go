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
func (k Keeper) GetAllLegacyPools(ctx sdk.Context) (list []types.LegacyPool) {
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

func (k Keeper) GetAllLegacyPools(ctx sdk.Context) []types.LegacyPool {
	var pools []types.LegacyPool
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer func(iterator sdk.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		var pool types.LegacyPool
		bytesValue := iterator.Value()
		err := k.cdc.Unmarshal(bytesValue, &pool)
		if err == nil {
			pools = append(pools, pool)
		}
	}
	return pools
}

func (k Keeper) DeleteLegacyPool(ctx sdk.Context,  poolId uint64) error {
	store := ctx.KVStore(k.storeKey)
	key := types.PoolKey(poolId)
	if !store.Has(key) {
		return types.ErrPositionDoesNotExist
	}
	store.Delete(key)
	return nil
}
