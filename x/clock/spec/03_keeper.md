<!--
order: 3
-->

# Keeper

## Contract Management

The `clock` module's keeper handles the management of contracts, including retrieving and setting parameters.

### SetParams

The `SetParams` function sets the specific parameters in the store.

```go
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
    if err := p.Validate(); err != nil {
        return err
    }

    store := ctx.KVStore(k.storeKey)
    bz := k.cdc.MustMarshal(&p)
    store.Set(types.ParamsKey, bz)

    return nil
}
```

### GetParams

The `GetParams` function retrieves the parameters from the store.

```go
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get(types.ParamsKey)
    if bz == nil {
        return p
    }

    k.cdc.MustUnmarshal(bz, &p)
    return p
}
```
