<!--
order: 6
-->

# Functions

## BeginBlocker

### BeginBlocker

The `BeginBlocker` function checks if the epoch has passed and executes necessary operations for each pool.

```go
func (k Keeper) BeginBlocker(ctx sdk.Context) {
    // check if epoch has passed then execute
    epochLength := k.GetEpochLength(ctx)
    epochPosition := k.GetEpochPosition(ctx, epochLength)

    // if epoch has not passed
    if epochPosition != 0 {
        return
    }

    // if epoch has passed
    entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
    if !found {
        ctx.Logger().Error(errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency).Error())
    }
    baseCurrency := entry.Denom
    baseCurrencyDecimal := entry.Decimals

    currentHeight := ctx.BlockHeight()
    pools := k.GetAllPools(ctx)
    for _, pool := range pools {
        ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
        if err != nil {
            ctx.Logger().Error(errorsmod.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
            continue
        }
        if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
            rate, err := k.BorrowInterestRateComputation(ctx, pool)
            if err != nil {
                ctx.Logger().Error(err.Error())
                continue
            }
            pool.BorrowInterestRate = rate
            pool.LastHeightBorrowInterestRateComputed = currentHeight
            err = k.UpdatePoolHealth(ctx, &pool)
            if err != nil {
                ctx.Logger().Error(err.Error())
            }
            err = k.UpdateFundingRate(ctx, &pool)
            if err != nil {
                ctx.Logger().Error(err.Error())
            }

            mtps, _, _ := k.GetMTPsForPool(ctx, pool.AmmPoolId, nil)
            for _, mtp := range mtps {
                err := BeginBlockerProcessMTP(ctx, k, mtp, pool, ammPool, baseCurrency, baseCurrencyDecimal)
                if err != nil {
                    ctx.Logger().Error(err.Error())
                    continue
                }
            }
            err = k.HandleFundingFeeDistribution(ctx, mtps, &pool, ammPool, baseCurrency)
            if err != nil {
                ctx.Logger().Error(err.Error())
            }
        }
        k.SetPool(ctx, pool)
    }
}
```

## Params

### GetParams

The `GetParams` function retrieves the current parameters of the `perpetual` module.

```go
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get(types.KeyPrefix(types.ParamsKey))
    if bz == nil {
        return params
    }

    k.cdc.MustUnmarshal(bz, &params)
    return params
}
```

### SetParams

The `SetParams` function sets the parameters of the `perpetual` module.

```go
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
    if err := params.Validate(); err != nil {
        return err
    }

    store := ctx.KVStore(k.storeKey)
    bz, err := k.cdc.Marshal(params)
    if err != nil {
        return err
    }
    store.Set(types.KeyPrefix(types.ParamsKey), bz)

    return nil
}
```

## MTP

### GetMTP

The `GetMTP` function retrieves an MTP by address and ID.

```go
func (k Keeper) GetMTP(ctx sdk.Context, mtpAddress string, id uint64) (types.MTP, error) {
    var mtp types.MTP
    key := types.GetMTPKey(mtpAddress, id)
    store := ctx.KVStore(k.storeKey)
    if !store.Has(key) {
        return mtp, types.ErrMTPDoesNotExist
    }
    bz := store.Get(key)
    k.cdc.MustUnmarshal(bz, &mtp)
    return mtp, nil
}
```

### SetMTP

The `SetMTP` function sets an MTP in the store.

```go
func (k Keeper) SetMTP(ctx sdk.Context, mtp *types.MTP) error {
    store := ctx.KVStore(k.storeKey)
    count := k.GetMTPCount(ctx)
    openCount := k.GetOpenMTPCount(ctx)

    if mtp.Id == 0 {
        // increment global id count
        count++
        mtp.Id = count
        k.SetMTPCount(ctx, count)
        // increment open mtp count
        openCount++
        k.SetOpenMTPCount(ctx, openCount)
    }

    if err := mtp.Validate(); err != nil {
        return err
    }
    key := types.GetMTPKey(mtp.Address, mtp.Id)
    store.Set(key, k.cdc.MustMarshal(mtp))
    return nil
}
```

## Pool

### GetPool

The `GetPool` function returns a pool by its ID.

```go
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))

    b := store.Get(types.PoolKey(poolId))
    if b == nil {
        return val, false
    }

    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

### SetPool

The `SetPool` function sets a specific pool in the store by its index.

```go
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolKeyPrefix))
    b := k.cdc.MustMarshal(&pool)
    store.Set(types.PoolKey(pool.AmmPoolId), b)
}
```

## Hooks

### SetHooks

The `SetHooks` function sets the perpetual hooks.

```go
func (k *Keeper) SetHooks(gh types.PerpetualHooks) *Keeper {
    if k.hooks != nil {
        panic("cannot set perpetual hooks twice")
    }

    k.hooks = gh

    return k
}
```

## Query

### GetPositions

The `GetPositions` function retrieves MTPs with pagination.

```go
func (k Keeper) GetPositions(goCtx context.Context, req *types.PositionsRequest) (*types.PositionsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    if req.Pagination != nil && req.Pagination.Limit > types.MaxPageLimit {
        return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
    }

    ctx := sdk.UnwrapSDKContext(goCtx)
    mtps, page, err := k.GetMTPs(ctx, req.Pagination)
    if err != nil {
        return nil, err
    }

    return &types.PositionsResponse{
        Mtps:       mtps,
        Pagination: page,
    }, nil
}
```

### GetPositionsByPool

The `GetPositionsByPool` function retrieves MTPs for a specific pool with pagination.

```go
func (k Keeper) GetPositionsByPool(goCtx context.Context, req *types.PositionsByPoolRequest) (*types.PositionsByPoolResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    ctx := sdk.UnwrapSDKContext(goCtx)
    mtps, pageRes, err := k.GetMTPsForPool(ctx, req.AmmPoolId, req.Pagination)
    if err != nil {
        return nil, err
    }

    return &types.PositionsByPoolResponse{
        Mtps:       mtps,
        Pagination: pageRes,
    }, nil
}
```

### GetPositionsForAddress

The `GetPositionsForAddress` function retrieves MTPs for a specific address with pagination.

```go
func (k Keeper) GetPositionsForAddress(goCtx context.Context, req *types.PositionsForAddressRequest) (*types.PositionsForAddressResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    addr, err := sdk.AccAddressFromBech32(req.Address)
    if err != nil {
        return nil, err
    }

    mtps, pageRes, err := k.GetMTPsForAddressWithPagination(sdk.UnwrapSDKContext(goCtx), addr, req.Pagination)
    if err != nil {
        return nil, err
    }

    return &types.PositionsForAddressResponse{Mtps: mtps, Pagination: pageRes}, nil
}
```

### GetStatus

The `GetStatus` function retrieves the open and lifetime MTP count.

```go
func (k Keeper) GetStatus(goCtx context.Context, req *types.StatusRequest) (*types.StatusResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    ctx := sdk.UnwrapSDKContext(goCtx)
    return &types.StatusResponse{
        OpenMtpCount:     k.GetOpenMTPCount(ctx

),
        LifetimeMtpCount: k.GetMTPCount(ctx),
    }, nil
}
```

### GetModuleBalance

The `GetModuleBalance` function retrieves the balance of the module account.

```go
func (k Keeper) GetModuleBalance(ctx sdk.Context, denom string) sdk.Coin {
    moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
    balance := k.bankKeeper.GetBalance(ctx, moduleAccount.GetAddress(), denom)
    return balance
}
```

### GetAllPools

The `GetAllPools` function retrieves all pools from the store.

```go
func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.PoolKeyPrefix))

    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var pool types.Pool
        k.cdc.MustUnmarshal(iterator.Value(), &pool)
        pools = append(pools, pool)
    }

    return pools
}
```

## Epoch

### GetEpochLength

The `GetEpochLength` function retrieves the length of an epoch.

```go
func (k Keeper) GetEpochLength(ctx sdk.Context) uint64 {
    return k.GetParams(ctx).EpochLength
}
```

### GetEpochPosition

The `GetEpochPosition` function calculates the current position within the epoch.

```go
func (k Keeper) GetEpochPosition(ctx sdk.Context, epochLength uint64) uint64 {
    return uint64(ctx.BlockHeight()) % epochLength
}
```

## MTP Handling

### BeginBlockerProcessMTP

The `BeginBlockerProcessMTP` function processes an MTP during the BeginBlocker.

```go
func BeginBlockerProcessMTP(
    ctx sdk.Context,
    k Keeper,
    mtp types.MTP,
    pool types.Pool,
    ammPool types.AmmPool,
    baseCurrency string,
    baseCurrencyDecimal uint8,
) error {
    // Implementation of the processing logic
    // ...

    return nil
}
```

### HandleFundingFeeDistribution

The `HandleFundingFeeDistribution` function handles the distribution of funding fees.

```go
func (k Keeper) HandleFundingFeeDistribution(
    ctx sdk.Context,
    mtps []types.MTP,
    pool *types.Pool,
    ammPool types.AmmPool,
    baseCurrency string,
) error {
    // Implementation of the fee distribution logic
    // ...

    return nil
}
```

## Interest Rate

### BorrowInterestRateComputation

The `BorrowInterestRateComputation` function computes the borrow interest rate for a pool.

```go
func (k Keeper) BorrowInterestRateComputation(ctx sdk.Context, pool types.Pool) (sdk.Dec, error) {
    // Implementation of the interest rate computation
    // ...

    return sdk.ZeroDec(), nil
}
```

## Pool Health

### UpdatePoolHealth

The `UpdatePoolHealth` function updates the health of a pool.

```go
func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) error {
    // Implementation of the pool health update logic
    // ...

    return nil
}
```

### UpdateFundingRate

The `UpdateFundingRate` function updates the funding rate for a pool.

```go
func (k Keeper) UpdateFundingRate(ctx sdk.Context, pool *types.Pool) error {
    // Implementation of the funding rate update logic
    // ...

    return nil
}
```
