# `incentive` module

## Overview

The Incentive module is designed to reward the ecosystem participants including liquidity providers (LPs), Elys stakers, and Eden committers.

Reward is distributed per epoch, `distribution` epoch, which is counted in number of blocks (`distribution_interval`).

There's Eden allocation epoch per day, based on tokenomics.
The source of rewards are from `Eden + Dex revenue (USDC) + Gas fees (XX,YY -> USDC)`

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Keeper](03_keeper.md)**
4. **[Endpoints](04_endpoints.md)**

## Reference

[Notion page](https://www.notion.so/Incentives-Module-Spec-bc6547edaf26472fa92c877740e2cd12)
