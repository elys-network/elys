<!--
order: 3
-->

# Keeper

## Rewards Distribution

The `Masterchef` module's keeper handles the distribution of LP rewards and external incentives. It ensures that rewards are properly calculated and distributed and that necessary adjustments to staking parameters are made regularly.

### EndBlocker

The `EndBlocker` function is invoked at the end of each block. It is responsible for processing LP rewards and external incentives distribution.

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
    k.ProcessLPRewardDistribution(ctx)
    k.ProcessExternalRewardsDistribution(ctx)
}
```

### ProcessLPRewardDistribution

The `ProcessLPRewardDistribution` function distributes rewards to liquidity providers based on collected fees and staking conditions.

```go
func (k Keeper) ProcessLPRewardDistribution(ctx sdk.Context) {
    k.ProcessUpdateIncentiveParams(ctx)
    err := k.UpdateLPRewards(ctx)
    if err != nil {
        ctx.Logger().Error("Failed to update LP rewards", "error", err)
    }
}
```

### ProcessExternalRewardsDistribution

The `ProcessExternalRewardsDistribution` function distributes external incentives to the specified pools within the defined block range.

```go
func (k Keeper) ProcessExternalRewardsDistribution(ctx sdk.Context) {
    // Logic for distributing external rewards
}
```

### UpdateLPRewards

The `UpdateLPRewards` function updates the rewards for liquidity providers by calculating the total rewards based on the collected fees and staking conditions.

```go
func (k Keeper) UpdateLPRewards(ctx sdk.Context) error {
    // Logic for updating LP rewards
}
```

### UpdateAccPerShare

The `UpdateAccPerShare` function updates the accumulated reward per share for a specific pool.

```go
func (k Keeper) UpdateAccPerShare(ctx sdk.Context, poolId uint64, rewardDenom string, amount sdk.Int) {
    // Logic for updating accumulated reward per share
}
```
