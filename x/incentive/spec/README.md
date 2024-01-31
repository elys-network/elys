# `incentive` module

## Overview

The Incentive module is designed to reward the ecosystem participants including liquidity providers (LPs), Elys stakers, and Eden committers.

Reward is distributed per epoch, `distribution` epoch, which is counted in number of blocks (`distribution_interval`).

There's Eden allocation epoch per day, based on tokenomics.
The source of rewards are from `Eden + Dex revenue (USDC) + Gas fees (XX,YY -> USDC)`

## Flow

### Staking

1. allocation of daily eden is based on tokenomics which is in staking allocation
2. capped allocation of daily eden is based on 30% Apr
3. distribution is done every set epoch

### LM rewards

1. allocation of daily eden in based on tokenmics module which is for LM alllocatioj
2. capped allocation of daily eden for 50% Apr
3. weights for different pools as different pools will be given different rewards
4. usdc stable coin pool is included here
5. distribution based on proxy tVL ( weighted TVL)
6. distribution every set epoch

### EdenBoost allocation

Eden boost is received at 100% Apr for staking elys and committing Eden.
Eden staking and Elys staking is exactly the same other than securing the chain and bonding period
They both get Eden and Eden boost.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Keeper](03_keeper.md)**
4. **[Endpoints](04_endpoints.md)**

## Reference

[Notion page](https://www.notion.so/Incentives-Module-Spec-bc6547edaf26472fa92c877740e2cd12)
