<!--
order: 3
-->

# Keeper

## Rewards Distribution

The `Masterchef` module's keeper handles the distribution of LP rewards and external incentives. It ensures that rewards are properly calculated and distributed and that necessary adjustments to staking parameters are made regularly.

At every block, iteration of LP pools' done and Eden rewards are calculated based on proxy TVL. Pool proxy TVL is expressed as `Pool TVL * Pool Multiplier`.

When rewards are distributed to a pool, `UpdateAccPerShare` is called and it is updating `PoolRewardInfo` object which includes pool accumulated reward per share.

Reward distribution's executed through `EndBlocker`. The `EndBlocker` function is called at the end of each block to perform necessary updates and maintenance for the `masterchef` module. It processes LP rewards and external incentives distribution.

### ProcessLPRewardDistribution

The `ProcessLPRewardDistribution` function distributes rewards to liquidity providers. It updates the incentive parameters and calculates the rewards based on the collected fees and staking conditions.

Note: There's constant reward pool address per liquidity pool. `GetLPRewardsPoolAddress` generates an address to be used in collecting DEX revenue from the specified pool.

Gas fees collected are converted into USDC before being distributed. DEX swap fees are collected in USDC, therefore masterchef module doesn't convert DEX fees to USDC.

### ProcessExternalRewardsDistribution

The `ProcessExternalRewardsDistribution` function distributes external incentives to the specified pools within the defined block range.

### UpdateLPRewards

The `UpdateLPRewards` function updates the rewards for liquidity providers by calculating the total rewards based on the collected fees and staking conditions.

### UpdateAccPerShare

The `UpdateAccPerShare` function updates the accumulated reward per share for a specific pool.

### CollectGasFees

The `CollectGasFees` function collects gas fees and allocates them to LPs and the protocol.

### CollectDEXRevenue

The `CollectDEXRevenue` function collects DEX revenue and distributes it to LPs and the protocol.

### CalculateProxyTVL

The `CalculateProxyTVL` function calculates the proxy total value locked (TVL) for the pools.

### UpdateAmmPoolAPR

The `UpdateAmmPoolAPR` function updates the APR for AMM pools.

## Rewards claim

Rewards claim is done per pool and reward denom.
`UserRewardInfo`'s `RewardPending` and `RewardDebt` fields are updated and appropriate amount of tokens are sent from module to user account.
