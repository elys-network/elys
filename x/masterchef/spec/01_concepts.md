<!--
order: 1
-->

# Concepts

`masterchef` module is to support rewards distribution to liquidity providers.

The terminology `masterchef` came from the common reward distribution mechanism in EVM chains, e.g. [sushiswap masterchef reward distribution mechanism](https://dev.to/heymarkkop/understanding-sushiswaps-masterchef-staking-rewards-1m6f).

Masterchef is distributing rewards per block. There's Eden allocation, based on tokenomics to liquidity providers.

The source of rewards are from `Eden` allocation, `DEX revenue` (USDC) and `Gas fees` (All gas fees are swapped to USDC on amm pools before distribution.)

## Flow

1. Allocation of eden in based on tokenmics module which is capped allocation of eden for 50% Apr (Note: multiplier is applied on APR cap, e.g. if multiplier is 1.5 cap's increased to 75% APR)
2. Every block, gas and swap fee rewards are distributed to liquidity providers
3. Reward is splitted based on proxy TVL which is `TVL * multiplier` (USDC stable pool's included as one of pools with lower multiplier)
