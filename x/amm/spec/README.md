# Elys AMM Module

## Contents

1. **[Concepts](01_concepts.md)**
2. **[Mechanism](02_mechanism.md)**
3. **[State](03_state.md)**
4. **[Keeper](04_keeper.md)**
5. **[Endpoints](05_endpoints.md)**
6. **[Risk Management](06_risk_management.md)**
7. **[Slippage Model](07_slippage.md)**
8. **[Swap Transactions Reordering](08_swap_txs_reordering.md)**

## References

Resources:

- [Elys Network Documentation](https://docs.elys.network)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [GitHub Repository for Elys Network](https://github.com/elys-network/elys)

## Overview

The Elys Network AMM module implements an Automated Market Maker (AMM) using various types of liquidity pools. It supports the following types of pools:

1. **AMM Pools**: Non-oracle pools with liquidity centered around a specific spot price, designed for assets with significant price variation.
2. **Oracle Pools**: Oracle pools with liquidity centered around an oracle price, designed for assets with stable prices.
