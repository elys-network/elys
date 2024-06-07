<!--
order: 3
-->

# Keeper

## Airdrop Management

The `tokenomics` module's keeper manages airdrop entries, ensuring accurate storage, retrieval, and removal of airdrop data.

### Set Airdrop

Sets an airdrop entry in the store.

```go
func (k Keeper) SetAirdrop(ctx sdk.Context, airdrop types.Airdrop) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))
    b := k.cdc.MustMarshal(&airdrop)
    store.Set(types.AirdropKey(airdrop.Intent), b)
}
```

### Get Airdrop

Retrieves an airdrop entry by intent.

```go
func (k Keeper) GetAirdrop(ctx sdk.Context, intent string) (val types.Airdrop, found bool) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))

    b := store.Get(types.AirdropKey(intent))
    if b == nil {
        return val, false
    }

    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

### Remove Airdrop

Removes an airdrop entry by intent.

```go
func (k Keeper) RemoveAirdrop(ctx sdk.Context, intent string) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))
    store.Delete(types.AirdropKey(intent))
}
```

### GetAllAirdrop

Retrieves all airdrop entries.

```go
func (k Keeper) GetAllAirdrop(ctx sdk.Context) (list []types.Airdrop) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AirdropKeyPrefix))
    iterator := sdk.KVStorePrefixIterator(store, []byte{})

    defer iterator.Close()

    for ; iterator.Valid(); iterator.Next() {
        var val types.Airdrop
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }

    return
}
```

## Inflation Management

### Set Genesis Inflation

Sets the genesis inflation parameters in the store.

```go
func (k Keeper) SetGenesisInflation(ctx sdk.Context, genesisInflation types.GenesisInflation) {
    store := ctx.KVStore(k.storeKey)
    b := k.cdc.MustMarshal(&genesisInflation)
    store.Set([]byte(types.GenesisInflationKey), b)
}
```

### Get Genesis Inflation

Retrieves the genesis inflation parameters.

```go
func (k Keeper) GetGenesisInflation(ctx sdk.Context) (val types.GenesisInflation, found bool) {
    store := ctx.KVStore(k.storeKey)

    b := store.Get([]byte(types.GenesisInflationKey))
    if b == nil {
        return val, false
    }

    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

### Remove Genesis Inflation

Removes the genesis inflation parameters from the store.

```go
func (k Keeper) RemoveGenesisInflation(ctx sdk.Context) {
    store := ctx.KVStore(k.storeKey)
    store.Delete([]byte(types.GenesisInflationKey))
}
```

## Query Handlers

### Query Airdrop Entries

```go
func (k Keeper) AirdropAll(goCtx context.Context, req *types.QueryAllAirdropRequest) (*

types.QueryAllAirdropResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    var airdrops []types.Airdrop
    ctx := sdk.UnwrapSDKContext(goCtx)

    store := ctx.KVStore(k.storeKey)
    airdropStore := prefix.NewStore(store, types.KeyPrefix(types.AirdropKeyPrefix))

    pageRes, err := query.Paginate(airdropStore, req.Pagination, func(key []byte, value []byte) error {
        var airdrop types.Airdrop
        if err := k.cdc.Unmarshal(value, &airdrop); err != nil {
            return err
        }

        airdrops = append(airdrops, airdrop)
        return nil
    })
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &types.QueryAllAirdropResponse{Airdrop: airdrops, Pagination: pageRes}, nil
}
```

### Query Genesis Inflation

```go
func (k Keeper) GenesisInflation(goCtx context.Context, req *types.QueryGetGenesisInflationRequest) (*types.QueryGetGenesisInflationResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)

    val, found := k.GetGenesisInflation(ctx)
    if !found {
        return nil, status.Error(codes.NotFound, "not found")
    }

    return &types.QueryGetGenesisInflationResponse{GenesisInflation: val}, nil
}
```
