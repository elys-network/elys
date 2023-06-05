# `incentive` module

## Overview

The Incentive module is designed to reward the ecosystem participants including liquidity providers (LPs), stakers, and Eden token holders.

- Fees are collected from amm, margin, gas fees
- Additional reward tokens come from commitment module by increasing their uncommitted token balances periodically. This module is integrated with the commitment module and updates accounting for both Eden tokens and Eden-Boost tokens.
- Distribution targets are LP token and delegators.
- Reward tokens are in Eden, Eden Boost, USDC (ELYS token is never rewarded - ELYS token within the reward pool swapped to USDC)
- Community pool management

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Keeper](03_keeper.md)**
4. **[Endpoints](04_endpoints.md)**

## Reference

[Notion page](https://www.notion.so/Incentives-Module-Spec-bc6547edaf26472fa92c877740e2cd12)
