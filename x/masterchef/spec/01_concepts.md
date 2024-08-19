<!--
order: 1
-->

# Concepts

The masterchef module is designed to facilitate rewards distribution to liquidity providers. The term masterchef is derived from a widely recognized reward distribution mechanism used in EVM chains, such as the [Sushiswap Masterchef reward distribution mechanism](https://dev.to/heymarkkop/understanding-sushiswaps-masterchef-staking-rewards-1m6f).

The masterchef module distributes rewards on a per-block basis. The sources of these rewards include:

- Eden Allocation: Governed by the tokenomics module.
- DEX Revenue (USDC): Generated from decentralized exchange operations.
- Gas Fees: All gas fees are swapped to USDC via AMM pools before distribution.

## Flow

1. Eden Allocation:
    - Eden is allocated based on the tokenomics module, with a capped APR set at 50%.
    - This cap can be adjusted using a multiplier; for example, a multiplier of 1.5 increases the APR cap to 75%.
2. Rewards Distribution:
    - Rewards from gas fees and swap fees are distributed to liquidity providers on a per-block basis.
    - All liquidity provider tokens are locked within the commitment module, where they contribute to liquidity and are used for reward calculations.
    - Rewards are calculated based on the committed LP tokens within the commitment module.
3. Reward Splitting:
    - Rewards are distributed based on the proxy TVL, calculated as TVL * multiplier.
    - The USDC stable pool is included in the TVL calculation, but it is assigned a lower multiplier compared to other pools.
