<!--
order: 5
-->

# Functions

## InitiateAccountedPool

The `InitiateAccountedPool` function initializes a new accounted pool based on an AMM pool.

```go
func (k Keeper) InitiateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
    poolId := ammPool.PoolId
    exists := k.PoolExists(ctx, poolId)
    if exists {
        return types.ErrPoolDoesNotExist
    }
    accountedPool := types.AccountedPool{
        PoolId:      poolId,
        TotalShares: ammPool.TotalShares,
        PoolAssets:  []ammtypes.PoolAsset{},
        TotalWeight: ammPool.TotalWeight,
    }
    for _, asset := range ammPool.PoolAssets {
        accountedPool.PoolAssets = append(accountedPool.PoolAssets, asset)
    }
    k.SetAccountedPool(ctx, accountedPool)
    return nil
}
```

## UpdateAccountedPool

The `UpdateAccountedPool` function updates an existing accounted pool based on AMM and Perpetual pool states.

```go
func (k Keeper) UpdateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) error {
    poolId := ammPool.PoolId
    exists := k.PoolExists(ctx, poolId)
    if !exists {
        return types.ErrPoolDoesNotExist
    }
    accountedPool, found := k.GetAccountedPool(ctx, poolId)
    if !found {
        return types.ErrPoolDoesNotExist
    }
    for i, asset := range accountedPool.PoolAssets {
        aBalance, err := perpetualtypes.GetAmmPoolBalance(ammPool, asset.Token.Denom)
        if err != nil {
            return err
        }
        mBalance, mLiabiltiies, _ := perpetualPool.GetPerpetualPoolBalances(asset.Token.Denom)
        accountedAmt := aBalance.Add(mBalance).Add(mLiabiltiies)
        accountedPool.PoolAssets[i].Token = sdk.NewCoin(asset.Token.Denom, accountedAmt)
    }
    k.SetAccountedPool(ctx, accountedPool)
    return nil
}
```

## Hooks

The `accountedpool` module includes several hooks to integrate with other modules such as AMM and Perpetual pools.

### AfterAmmPoolCreated

The `AfterAmmPoolCreated` hook is called after a new AMM pool is created

, initiating the corresponding accounted pool.

```go
func (k Keeper) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool) {
    k.InitiateAccountedPool(ctx, ammPool)
}
```

### AfterAmmJoinPool

The `AfterAmmJoinPool` hook is called after a join pool event, updating the corresponding accounted pool.

```go
func (k Keeper) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

### AfterAmmExitPool

The `AfterAmmExitPool` hook is called after an exit pool event, updating the corresponding accounted pool.

```go
func (k Keeper) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

### AfterAmmSwap

The `AfterAmmSwap` hook is called after a swap event, updating the corresponding accounted pool.

```go
func (k Keeper) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

### AfterPerpetualPositionOpen

The `AfterPerpetualPositionOpen` hook is called after a perpetual position is opened, updating the corresponding accounted pool.

```go
func (k Keeper) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

### AfterPerpetualPositionModified

The `AfterPerpetualPositionModified` hook is called after a perpetual position is modified, updating the corresponding accounted pool.

```go
func (k Keeper) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

### AfterPerpetualPositionClosed

The `AfterPerpetualPositionClosed` hook is called after a perpetual position is closed, updating the corresponding accounted pool.

```go
func (k Keeper) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) {
    k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}
```

### PerpetualHooks Wrapper

The `PerpetualHooks` wrapper struct for the `accountedpool` keeper implements the `PerpetualHooks` interface.

```go
type PerpetualHooks struct {
    k Keeper
}

var _ perpetualtypes.PerpetualHooks = PerpetualHooks{}

func (k Keeper) PerpetualHooks() PerpetualHooks {
    return PerpetualHooks{k}
}
```
