<!--
order: 5
-->

# Functions

## CreateAirdrop

The `CreateAirdrop` function creates a new airdrop entry.

```go
func (k msgServer) CreateAirdrop(goCtx context.Context, msg *types.MsgCreateAirdrop) (*types.MsgCreateAirdropResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value already exists
    _, found := k.GetAirdrop(ctx, msg.Intent)
    if found {
        return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
    }

    airdrop := types.Airdrop{
        Authority: msg.Authority,
        Intent:    msg.Intent,
        Amount:    msg.Amount,
        Expiry:    msg.Expiry,
    }

    k.SetAirdrop(ctx, airdrop)
    return &types.MsgCreateAirdropResponse{}, nil
}
```

## UpdateAirdrop

The `UpdateAirdrop` function updates an existing airdrop entry.

```go
func (k msgServer) UpdateAirdrop(goCtx context.Context, msg *types.MsgUpdateAirdrop) (*types.MsgUpdateAirdropResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, found := k.GetAirdrop(ctx, msg.Intent)
    if !found {
        return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Authority != valFound.Authority {
        return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    airdrop := types.Airdrop{
        Authority: msg.Authority,
        Intent:    msg.Intent,
        Amount:    msg.Amount,
        Expiry:    msg.Expiry,
    }

    k.SetAirdrop(ctx, airdrop)
    return &types.MsgUpdateAirdropResponse{}, nil
}
```

## DeleteAirdrop

The `DeleteAirdrop` function deletes an existing airdrop entry.

```go
func (k msgServer) DeleteAirdrop(goCtx context.Context, msg *types.MsgDeleteAirdrop) (*types.MsgDeleteAirdropResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, found := k.GetAirdrop(ctx, msg.Intent)
    if !found {
        return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Authority != valFound.Authority {
        return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    k.RemoveAirdrop(ctx, msg.Intent)
    return &types.MsgDeleteAirdropResponse{}, nil
}
```

## ClaimAirdrop

The `ClaimAirdrop` function allows a user to claim an airdrop.

```go
func (k msgServer) ClaimAirdrop(goCtx context.Context, msg *types.MsgClaimAirdrop) (*types.MsgClaimAirdropResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    airdrop, found := k.GetAirdrop(ctx, msg.Sender)
    if !found {
        return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Sender != airdrop.Authority {
        return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    if ctx.BlockTime().Unix() > int64(airdrop.Expiry) {
        return nil, types.ErrAirdropExpired
    }

    // Add commitments
    commitments := k.commitmentKeeper.GetCommitments(ctx, msg.Sender)
    commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, math.NewInt(int64(airdrop.Amount))))
    k.commitmentKeeper.SetCommitments(ctx, commitments)

    k.RemoveAirdrop(ctx, msg.Sender)
    return &types.MsgClaimAirdropResponse{}, nil
}
```

## UpdateGenesisInflation

The `UpdateGenesisInflation` function updates the genesis inflation parameters.

```go
func (k msgServer) UpdateGenesisInflation(goCtx context.Context, msg *types.MsgUpdateGenesisInflation) (*types.MsgUpdateGenesisInflationResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    genesisInflation := types.GenesisInflation{
        Authority:             msg.Authority,
        Inflation:             msg.Inflation,
        SeedVesting:           msg.SeedVesting,
        StrategicSalesVesting: msg.StrategicSalesVesting,
    }

    k.SetGenesisInflation(ctx, genesisInflation)

    return &types.MsgUpdateGenesisInflationResponse{}, nil
}
```

## CreateTimeBasedInflation

The `CreateTimeBasedInflation` function creates a new time-based inflation entry.

```go
func (k msgServer) CreateTimeBasedInflation(goCtx context.Context, msg *types.MsgCreateTimeBasedInflation) (*types.MsgCreateTimeBasedInflationResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value already exists
    _, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
    if found {
        return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
    }

    timeBasedInflation := types.TimeBasedInflation{
        Authority:        msg.Authority,
        StartBlockHeight: msg.StartBlockHeight,
        EndBlockHeight:   msg.EndBlockHeight,
        Description:      msg.Description,
        Inflation:        msg.Inflation,
    }

    k.SetTimeBasedInflation(ctx, timeBasedInflation)
    return &types.MsgCreateTimeBasedInflationResponse{}, nil
}
```

## UpdateTimeBasedInflation

The `UpdateTimeBasedInflation` function updates an existing time-based inflation entry

.

```go
func (k msgServer) UpdateTimeBasedInflation(goCtx context.Context, msg *types.MsgUpdateTimeBasedInflation) (*types.MsgUpdateTimeBasedInflationResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
    if !found {
        return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Authority != valFound.Authority {
        return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    timeBasedInflation := types.TimeBasedInflation{
        Authority:        msg.Authority,
        StartBlockHeight: msg.StartBlockHeight,
        EndBlockHeight:   msg.EndBlockHeight,
        Description:      msg.Description,
        Inflation:        msg.Inflation,
    }

    k.SetTimeBasedInflation(ctx, timeBasedInflation)

    return &types.MsgUpdateTimeBasedInflationResponse{}, nil
}
```

## DeleteTimeBasedInflation

The `DeleteTimeBasedInflation` function deletes an existing time-based inflation entry.

```go
func (k msgServer) DeleteTimeBasedInflation(goCtx context.Context, msg *types.MsgDeleteTimeBasedInflation) (*types.MsgDeleteTimeBasedInflationResponse, error) {
    if k.authority != msg.Authority {
        return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
    }

    ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, found := k.GetTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
    if !found {
        return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg authority is the same as the current owner
    if msg.Authority != valFound.Authority {
        return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    k.RemoveTimeBasedInflation(ctx, msg.StartBlockHeight, msg.EndBlockHeight)
    return &types.MsgDeleteTimeBasedInflationResponse{}, nil
}
```

## QueryHandlers

### Query All Airdrop Entries

The `AirdropAll` function returns a list of all airdrop entries.

```go
func (k Keeper) AirdropAll(goCtx context.Context, req *types.QueryAllAirdropRequest) (*types.QueryAllAirdropResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    var airdrops []types.Airdrop
    ctx := sdk.UnwrapSDKContext(goCtx)

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
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

### Query Specific Airdrop Entry

The `Airdrop` function returns the details of a specific airdrop entry based on the intent.

```go
func (k Keeper) Airdrop(goCtx context.Context, req *types.QueryGetAirdropRequest) (*types.QueryGetAirdropResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)

    val, found := k.GetAirdrop(ctx, req.Intent)
    if !found {
        return nil, status.Error(codes.NotFound, "not found")
    }

    return &types.QueryGetAirdropResponse{Airdrop: val}, nil
}
```

### Query Genesis Inflation Parameters

The `GenesisInflation` function returns the genesis inflation parameters.

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

### Query All Time-Based Inflation Entries

The `TimeBasedInflationAll` function returns a list of all time-based inflation entries.

```go
func (k Keeper) TimeBasedInflationAll(goCtx context.Context, req *types.QueryAllTimeBasedInflationRequest) (*types.QueryAllTimeBasedInflationResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    var timeBasedInflations []types.TimeBasedInflation
    ctx := sdk.UnwrapSDKContext(goCtx)

    store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    timeBasedInflationStore := prefix.NewStore(store, types.KeyPrefix(types.TimeBasedInflationKeyPrefix))

    pageRes, err := query.Paginate(timeBasedInflationStore, req.Pagination, func(key []byte, value []byte) error {
        var timeBasedInflation types.TimeBasedInflation
        if err := k.cdc.Unmarshal(value, &timeBasedInflation); err != nil {
            return err
        }

        timeBasedInflations = append(timeBasedInflations, timeBasedInflation)
        return nil
    })
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    return &types.QueryAllTimeBasedInflationResponse{TimeBasedInflation: timeBasedInflations, Pagination: pageRes}, nil
}
```

### Query Specific Time-Based Inflation Entry

The `TimeBasedInflation` function returns the details of a specific time-based inflation entry based on start and end block height.

```go
func (k Keeper) TimeBasedInflation(goCtx context.Context, req *types.QueryGetTimeBasedInflationRequest) (*types.QueryGetTimeBasedInflationResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    ctx := sdk.UnwrapSDKContext(goCtx)

    val, found := k.GetTimeBasedInflation(
        ctx,
        req.StartBlockHeight,
        req.EndBlockHeight,
    )
    if !found {
        return nil, status.Error(codes.NotFound, "not found")
    }

    return &types.QueryGetTimeBasedInflationResponse{TimeBasedInflation: val}, nil
}
```
