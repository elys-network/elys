<!--
order: 3
-->

# Keeper

## Accounted Pool Management

The `accountedpool` module's keeper handles the management of accounted pools, including creation, updates, and state retrieval.

### SetAccountedPool

The `SetAccountedPool` function sets a specific accounted pool in the store.

```go
func (k Keeper) SetAccountedPool(ctx sdk.Context, accountedPool types.AccountedPool) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
    b := k.cdc.MustMarshal(&accountedPool)
    store.Set(types.AccountedPoolKey(accountedPool.PoolId), b)
}
```

### GetAccountedPool

The `GetAccountedPool` function retrieves an accounted pool from the store by its ID.

```go
func (k Keeper) GetAccountedPool(ctx sdk.Context, PoolId uint64) (val types.AccountedPool, found bool) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
    b := store.Get(types.AccountedPoolKey(PoolId))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

### RemoveAccountedPool

The `RemoveAccountedPool` function removes an accounted pool from the store.

```go
func (k Keeper) RemoveAccountedPool(ctx sdk.Context, poolId uint64) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
    store.Delete(types.AccountedPoolKey(poolId))
}
```

### GetAllAccountedPool

The `GetAllAccountedPool` function retrieves all accounted pools from the store.

```go
func (k Keeper) GetAllAccountedPool(ctx sdk.Context) (list []types.AccountedPool) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
    iterator := storetypes.KVStorePrefixIterator(store, []byte{})
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var val types.AccountedPool
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }
    return
}
```
