# Elys AMM Module

## Contents

1. **[Concepts](01_concepts.md)**
2. **[Mechanism](02_mechanism.md)**
3. **[Oracle Pools](03_oracle_pools.md)**
4. **[State](04_state.md)**
5. **[Keeper](05_keeper.md)**
6. **[Endpoints](06_endpoints.md)**
7. **[Risk Management](07_risk_management.md)**
8. **[Slippage Model](08_slippage.md)**
9. **[Swap Transactions Reordering](09_swap_txs_reordering.md)**

## References

Resources:

- [Elys Network Documentation](https://docs.elys.network)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [GitHub Repository for Elys Network](https://github.com/elys-network/elys)

## Overview

The Elys Network AMM module implements an Automated Market Maker (AMM) using various types of liquidity pools. It supports the following types of pools:

1. **AMM Pools**: Non-oracle pools with liquidity centered around a specific spot price, designed for assets with significant price variation.
2. **Oracle Pools**: Oracle pools with liquidity centered around an oracle price, designed for assets with stable prices.
