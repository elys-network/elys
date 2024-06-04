# masterchef module

`masterchef` module is to support rewards distribution to liquidity providers.

## Overview

Masterchef is distributing rewards per block.

There's Eden allocation, based on tokenomics to liquidity providers.
The source of rewards are from `Eden + Dex revenue (USDC) + Gas fees (XX,YY -> USDC)`

## Flow

1. Allocation of eden in based on tokenmics module which is capped allocation of eden for 50% Apr (Note: multiplier is applied on APR cap, e.g. if multiplier is 1.5 cap's increased to 75% APR)
2. Every block, gas and swap fee rewards are distributed to liquidity providers
3. Reward is splitted based on proxy TVL which is `TVL * multiplier` (USDC stable pool's included as one of pools with lower multiplier)

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Keeper](03_keeper.md)**
4. **[Endpoints](04_endpoints.md)**
