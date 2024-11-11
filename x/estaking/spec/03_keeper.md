<!--
order: 3
-->

# Keeper

## Rewards Distribution

The `estaking` module's keeper handles rewards distribution and updates staking parameters. It ensures that rewards are properly calculated, distributed, and that necessary adjustments to staking parameters are made regularly.

### EndBlocker

The `EndBlocker` function is invoked at the end of each block. It is responsible for processing the distribution of rewards and burning EdenB tokens if there has been a reduction in Elys staking.

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
    k.ProcessRewardsDistribution(ctx)
    k.BurnEdenBIfElysStakingReduced(ctx)
}
```

### Taking Delegation Snapshot

The `TakeDelegationSnapshot` function captures the current state of a delegator's staked amount. This snapshot includes calculating the total delegation amount and storing it as an `ElysStaked` object.

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

### Burning EdenB Tokens

The `BurnEdenBIfElysStakingReduced` function burns EdenB tokens if the Elys staking has decreased. It checks for addresses where staking changes have occurred and performs necessary actions to burn tokens and update snapshots.

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

### Processing Rewards Distribution

The `ProcessRewardsDistribution` function is responsible for distributing rewards to stakers. It updates incentive parameters and calculates the rewards to be distributed based on collected fees and staking conditions.

```go
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) {
    k.ProcessUpdateIncentiveParams(ctx)
    err := k.UpdateStakersRewards(ctx)
    if err != nil {
        ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
    }
}
```

### Updating Stakers Rewards

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

    params.DexRewardsStakers.NumBlocks = sdkmath.OneInt()
    params.DexRewardsStakers.Amount = dexRevenueStakersAmount
    k.SetParams(ctx, params)

    coins := sdk.NewCoins(
        sdk.NewCoin(ptypes.Eden, stakersEdenAmount),
        sdk.NewCoin(ptypes.EdenB, stakersEdenBAmount),
    )
    return k.commKeeper.MintCoins(ctx, authtypes.FeeCollectorName, coins.Sort())
}
```
