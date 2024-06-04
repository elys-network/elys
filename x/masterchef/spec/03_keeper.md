<!--
order: 3
-->

# Keeper

## Reward pool address for LPs

There's constant reward pool address per liquidity pool.

`GetLPRewardsPoolAddress` generates an address to be used in collecting DEX revenue from the specified pool.

## Gas fees to USDC conversion

Gas fees collected are converted into USDC before being distributed. Other DEX fees will be collected in USDC. Meaning masterchef module doesn't need to work on converting DEX fees collected.

## Rewards distribution logic for LPs

We need to iterate LP pools and then calculate the rewards amount for LPs of the specified pool.
We use proxy TVL rather than pure TVL for calculating each pool share. Given the Eden amount per block, we use multiplier of each pool.

Each pool has unique LP token and it is committed to commitment module.
When rewards distributed to a pool, `UpdateAccPerShare` is called and it is updating `PoolRewardInfo` object which includes pool accumulated reward per share.

## Rewards claim

Rewards claim is done per pool and reward denom.
`UserRewardInfo`'s `RewardPending` and `RewardDebt` fields are updated and appropriate amount of tokens are sent from module to user account.
