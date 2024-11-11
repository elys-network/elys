<!--
order: 3
-->

# Keeper

## Parameter Management

The `parameter` module's keeper handles the management of configuration parameters, including setting, updating, and retrieving their values.

### GetParams

The `GetParams` function retrieves the current parameters from the store.

```go
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    b := store.Get([]byte(types.ParamsKey))
    if b == nil {
        return
    }
    k.cdc.MustUnmarshal(b, &params)
    return
}
```

### SetParams

The `SetParams` function sets the parameters in the store.

```go
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    b := k.cdc.MustMarshal(&params)
    store.Set([]byte(types.ParamsKey), b)
}
```

### GetLegacyParams

The `GetLegacyParams` function retrieves the legacy parameters from the store.

```go
func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    b := store.Get([]byte(types.ParamsKey))
    if b == nil {
        return
    }
    k.cdc.MustUnmarshal(b, &params)
    return
}
```
