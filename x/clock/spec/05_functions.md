<!--
order: 5
-->

# Functions

## EndBlocker

The `EndBlocker` function executes specified contracts at the end of each block.

```go
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
    defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

    message := []byte(types.EndBlockSudoMessage)
    p := k.GetParams(ctx)
    errorExecs := make([]string, 0)

    for idx, addr := range p.ContractAddresses {
        contract, err := sdk.AccAddressFromBech32(addr)
        if err != nil {
            errorExecs[idx] = addr
            continue
        }

        childCtx := ctx.WithGasMeter(sdk.NewGasMeter(p.ContractGasLimit))
        _, err = k.GetContractKeeper().Sudo(childCtx, contract, message)
        if err != nil {
            errorExecs = append(errorExecs, err.Error())
            continue
        }
    }

    if len(errorExecs) > 0 {
        log.Printf("[x/clock] Execute Errors: %v", errorExecs)
    }
}
```

## SetParams

The `SetParams` function sets new parameters for the module.

```go
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
    if err := p.Validate(); err != nil {
        return err
    }

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    bz := k.cdc.MustMarshal(&p)
    store.Set(types.ParamsKey, bz)

    return nil
}
```

## GetParams

The `GetParams` function retrieves the current parameters for the module.

```go
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    bz := store.Get(types.ParamsKey)
    if bz == nil {
        return p
    }

    k.cdc.MustUnmarshal(bz, &p)
    return p
}
```

## Hooks

The `clock` module does not include any hooks for integration with other modules.
