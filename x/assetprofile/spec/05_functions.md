<!--
order: 5
-->

# Functions

## Setting an Asset Entry

The `SetEntry` function stores a new asset entry or updates an existing one.

```go
func (k Keeper) SetEntry(ctx sdk.Context, entry types.Entry) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    b := k.cdc.MustMarshal(&entry)
    store.Set(types.EntryKey(entry.BaseDenom), b)
}
```

## Getting an Asset Entry

The `GetEntry` function retrieves an asset entry based on its base denomination.

```go
func (k Keeper) GetEntry(ctx sdk.Context, baseDenom string) (val types.Entry, found bool) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    b := store.Get(types.EntryKey(baseDenom))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

## Removing an Asset Entry

The `RemoveEntry` function deletes an asset entry from the store.

```go
func (k Keeper) RemoveEntry(ctx sdk.Context, baseDenom string) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    store.Delete(types.EntryKey(baseDenom))
}
```

## Retrieving All Asset Entries

The `GetAllEntry` function retrieves all asset entries from the store.

```go
func (k Keeper) GetAllEntry(ctx sdk.Context) (list []types.Entry) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    iterator := storetypes.KVStorePrefixIterator(store, []byte{})

    defer iterator.Close()

    for ; iterator.Valid(); iterator.Next() {
        var val types.Entry
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }

    return
}
```

## Creating an Asset Entry

The `CreateEntry` function handles the creation of a new asset entry.

```go
func (k msgServer) CreateEntry(goCtx context.Context, msg *types.MsgCreateEntry) (*types.MsgCreateEntryResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the entry already exists
    _, isFound := k.GetEntry(ctx, msg.BaseDenom)
    if isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "entry already set")
    }

    // check the validity of ibc denom & channel
    hash, err := ibctransfertypes.ParseHexHash(strings.TrimPrefix(msg.Denom, "ibc/"))
    if err == nil && k.transferKeeper != nil {
        denomTrace, ok := k.transferKeeper.GetDenomTrace(ctx, hash)
        if !ok {
            return nil, types.ErrNotValidIbcDenom
        }
        if !strings.Contains(denomTrace.Path, msg.IbcChannelId) {
            return nil, types.ErrChannelIdAndDenomHashMismatch
        }
    }

    entry := types.Entry{
        Authority:                msg.Authority,
        BaseDenom:                msg.BaseDenom,
        Decimals:                 msg.Decimals,
        Denom:                    msg.Denom,
        Path:                     msg.Path,
        IbcChannelId:             msg.IbcChannelId,
        IbcCounterpartyChannelId: msg.IbcCounterpartyChannelId,
        DisplayName:              msg.DisplayName,
        DisplaySymbol:            msg.DisplaySymbol,
        Network:                  msg.Network,
        Address:                  msg.Address,
        ExternalSymbol:           msg.ExternalSymbol,
        TransferLimit:            msg.TransferLimit,
        Permissions:              msg.Permissions,
        UnitDenom:                msg.UnitDenom,
        IbcCounterpartyDenom:     msg.IbcCounterpartyDenom,
        IbcCounterpartyChainId:   msg.IbcCounterpartyChainId,
        CommitEnabled:            msg.CommitEnabled,
        WithdrawEnabled:          msg.WithdrawEnabled,
    }

    k.SetEntry(ctx, entry)
    return &types.MsgCreateEntryResponse{}, nil
}
```

## Updating an Asset Entry

The `UpdateEntry` function handles the updating of an existing asset entry.

```go
func (k msgServer) UpdateEntry(goCtx context.Context, msg *types.MsgUpdateEntry) (*types.MsgUpdateEntryResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    entry, isFound := k.GetEntry(ctx, msg.BaseDenom)
    if !isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "entry not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Authority != entry.Authority {
        return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    // check the validity of ibc denom & channel
    hash, err := ibctransfertypes.ParseHexHash(strings.TrimPrefix(msg.Denom, "ibc/"))
    if err == nil && k.transferKeeper != nil {
        denomTrace, ok := k.transferKeeper.GetDenomTrace(ctx, hash)
        if !ok {
            return nil, types.ErrNotValidIbcDenom
        }
        if !strings.Contains(denomTrace.Path, msg.IbcChannelId) {
            return nil, types.ErrChannelIdAndDenomHashMismatch
        }
    }

    entry = types.Entry{
        Authority:                msg.Authority,
        BaseDenom:                msg.BaseDenom,
        Decimals:                 msg.Decimals,
        Denom:                    msg.Denom,
        Path:                     msg.Path,
        IbcChannelId:             msg.IbcChannelId,
        IbcCounterpartyChannelId: msg.IbcCounterpartyChannelId,
        DisplayName:              msg.DisplayName,
        DisplaySymbol:            msg.DisplaySymbol,
        Network:                  msg.Network,
        Address:                  msg.Address,
        ExternalSymbol:           msg.ExternalSymbol,
        TransferLimit:            msg.TransferLimit,
        Permissions:              msg.Permissions,
        UnitDenom:                msg.UnitDenom,
        IbcCounterpartyDenom:     msg.IbcCounterpartyDenom,
        IbcCounterpartyChainId:   msg.IbcCounterpartyChainId,
        CommitEnabled:            msg.CommitEnabled,
        WithdrawEnabled:          msg.WithdrawEnabled

,
    }

    k.SetEntry(ctx, entry)

    return &types.MsgUpdateEntryResponse{}, nil
}
```

## Deleting an Asset Entry

The `DeleteEntry` function handles the deletion of an existing asset entry.

```go
func (k msgServer) DeleteEntry(goCtx context.Context, msg *types.MsgDeleteEntry) (*types.MsgDeleteEntryResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    entry, isFound := k.GetEntry(ctx, msg.BaseDenom)
    if !isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "entry not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Authority != entry.Authority {
        return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    k.RemoveEntry(ctx, msg.BaseDenom)

    return &types.MsgDeleteEntryResponse{}, nil
}
```

## Querying Parameters

The `Params` function handles querying the parameters of the `assetprofile` module.

```go
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)

    return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
```

## Querying an Asset Entry

The `Entry` function handles querying an asset entry by its base denomination.

```go
func (k Keeper) Entry(goCtx context.Context, req *types.QueryEntryRequest) (*types.QueryEntryResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)

    val, found := k.GetEntry(ctx, req.BaseDenom)
    if !found {
        return nil, status.Error(codes.NotFound, "not found")
    }

    return &types.QueryEntryResponse{Entry: val}, nil
}
```

## Querying an Asset Entry by Denomination

The `EntryByDenom` function handles querying an asset entry by its denomination.

```go
func (k Keeper) EntryByDenom(goCtx context.Context, req *types.QueryEntryByDenomRequest) (*types.QueryEntryByDenomResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)

    val, found := k.GetEntryByDenom(ctx, req.Denom)
    if !found {
        return nil, status.Error(codes.NotFound, "not found")
    }

    return &types.QueryEntryByDenomResponse{Entry: val}, nil
}
```

## Querying All Asset Entries

The `EntryAll` function handles querying all asset entries.

```go
func (k Keeper) EntryAll(goCtx context.Context, req *types.QueryAllEntryRequest) (*types.QueryAllEntryResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    var entries []types.Entry
    ctx := sdk.UnwrapSDKContext(goCtx)

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    entryStore := prefix.NewStore(store, types.KeyPrefix(types.EntryKeyPrefix))

    pageRes, err := query.Paginate(entryStore, req.Pagination, func(key []byte, value []byte) error {
        var entry types.Entry
        if err := k.cdc.Unmarshal(value, &entry); err != nil {
            return err
        }

        entries = append(entries, entry)
        return nil
    })
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &types.QueryAllEntryResponse{Entry: entries, Pagination: pageRes}, nil
}
```
