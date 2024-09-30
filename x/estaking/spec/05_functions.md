<!--
order: 5
-->

# Functions

## EndBlocker

The `EndBlocker` function is called at the end of each block to perform necessary updates and maintenance for the `estaking` module. It processes rewards distribution and handles the burning of EdenB tokens if the Elys staking has reduced.

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
    k.ProcessRewardsDistribution(ctx)
    k.BurnEdenBIfElysStakingReduced(ctx)
}
```

### ProcessRewardsDistribution

The `ProcessRewardsDistribution` function is responsible for distributing rewards to stakers. It updates the incentive parameters and staker rewards based on the collected fees and staking conditions.

```go
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) {
    k.ProcessUpdateIncentiveParams(ctx)
    err := k.UpdateStakersRewards(ctx)
    if err != nil {
        ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
    }
}
```

### BurnEdenBIfElysStakingReduced

The `BurnEdenBIfElysStakingReduced` function burns EdenB tokens if the Elys staking has reduced. It checks for addresses where staking has changed and takes appropriate action to burn tokens and update staking snapshots.

```go
func (k Keeper) BurnEdenBIfElysStakingReduced(ctx sdk.Context) {
    addrs := k.GetAllElysStakeChange(ctx)
    for _, delAddr := range addrs {
        k.BurnEdenBFromElysUnstaking(ctx, delAddr)
        k.TakeDelegationSnapshot(ctx, delAddr.String())
        k.RemoveElysStakeChange(ctx, delAddr)
    }
}
```

## TakeDelegationSnapshot

The `TakeDelegationSnapshot` function captures the current state of a delegator's staked amount. It calculates the delegation amount and records it.

```go
func (k Keeper) TakeDelegationSnapshot(ctx sdk.Context, addr string) {
    delAmount := k.CalcDelegationAmount(ctx, addr)
    elysStaked := types.ElysStaked{
        Address: addr,
        Amount:  delAmount,
    }
    k.SetElysStaked(ctx, elysStaked)
}
```

## UpdateStakersRewards

The `UpdateStakersRewards` function updates the rewards for stakers. It calculates the total rewards based on the collected fees and staking conditions, and then mints the appropriate amount of reward tokens.

```go
func (k Keeper) UpdateStakersRewards(ctx sdk.Context) error {
    baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
    if !found {
        return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
    }

    feeCollectorAddr := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
    totalFeesCollected := k.commKeeper.GetAllBalances(ctx, feeCollectorAddr)
    gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(totalFeesCollected...)
    dexRevenueStakersAmount := gasFeeCollectedDec.AmountOf(baseCurrency)

    params := k.GetParams(ctx)
    stakeIncentive := params.StakeIncentives
    totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear

    edenAmountPerYear := sdkmath.ZeroInt()
    if stakeIncentive != nil && stakeIncentive.EdenAmountPerYear.IsPositive() {
        edenAmountPerYear = stakeIncentive.EdenAmountPerYear
    }
    stakersEdenAmount := edenAmountPerYear.Quo(math.NewInt(totalBlocksPerYear))

    totalElysEdenEdenBStake := k.TotalBondedTokens(ctx)

    stakersMaxEdenAmount := params.MaxEdenRewardAprStakers.
        MulInt(totalElysEdenEdenBStake).
        QuoInt64(totalBlocksPerYear)

    stakersEdenAmount = sdk.MinInt(stakersEdenAmount, stakersMaxEdenAmount.TruncateInt())

    stakersEdenBAmount := sdkmath.LegacyNewDecFromInt(totalElysEdenEdenBStake).
        Mul(params.EdenBoostApr).
        QuoInt64(totalBlocksPerYear).
        RoundInt()

    params.DexRewardsStakers.NumBlocks = sdk.OneInt()
    params.DexRewardsStakers.Amount = dexRevenueStakersAmount
    k.SetParams(ctx, params)

    coins := sdk.NewCoins(
        sdk.NewCoin(ptypes.Eden, stakersEdenAmount),
        sdk.NewCoin(ptypes.EdenB, stakersEdenBAmount),
    )
    return k.commKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins.Sort())
}
```
