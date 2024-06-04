<!--
order: 1
-->

# Concepts

`masterchef` module is to support rewards distribution to liquidity providers.

The terminology `masterchef` came from the common reward distribution mechanism in EVM chains, e.g. [sushiswap masterchef reward distribution mechanism](https://dev.to/heymarkkop/understanding-sushiswaps-masterchef-staking-rewards-1m6f).

Masterchef is distributing rewards per block. There's Eden allocation, based on tokenomics to liquidity providers.

The source of rewards are from `Eden` allocation, `DEX revenue` (USDC) and `Gas fees` (All gas fees are swapped to USDC on amm pools before distribution.)
