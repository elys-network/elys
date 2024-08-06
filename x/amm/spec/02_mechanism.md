<!--
order: 2
-->

# Mechanism

The Elys Network AMM module implements an Automated Market Maker (AMM) using various types of liquidity pools. It supports the following types of pools:

1. **AMM Pools**: Non-oracle pools with liquidity centered around a specific spot price, designed for assets with significant price variation.
2. **Oracle Pools**: Oracle pools with liquidity centered around an oracle price, designed for assets with stable prices.

This section explains the fundamental mechanism of the AMM module and provides an overview of how the module's code is structured to support both types of pools.

The section **[Oracle Pools](03_oracle_pools.md)** provides additional details on the mechanism of Oracle Pools.

## Pool

### Creation of Pool

When a new pool is created, a fixed amount of 100 share tokens is minted and sent to the pool creator's account. The pool share token denomination follows the format `amm/pool/{poolID}`. By default, pool assets are sorted in alphabetical order. Pool creation requires a minimum of 2 and a maximum of 8 asset denominations.

A `PoolCreationFee` must be paid to create a pool, which helps prevent the creation of unnecessary or malicious pools.

### Joining Pool

Users can join a pool without swapping assets by using the `JoinPool` function. In this case, they provide the maximum amount of tokens (`TokenInMaxs`) they are willing to deposit. This parameter must include all the asset denominations in the pool, or none at all, otherwise, the transaction will be aborted. If no tokens are specified, the user's balance is used as the constraint.

The front end calculates the number of share tokens the user is eligible for at the moment the transaction is sent. The exact calculation of tokens required for the designated share is performed during transaction processing, ensuring the deposit does not exceed the maximum specified. Once validated, the appropriate number of pool share tokens are minted and sent to the user's account.

#### Existing Join Types

- `JoinPool`

### Exiting Pool

To exit a pool, users specify the minimum amount of tokens they are willing to receive in return for their pool shares. Unlike joining, exiting a pool incurs an exit fee set as a pool parameter. The user's share tokens are burned in the process. Exiting with a single asset is also possible.

Exiting is only allowed if the user will leave a positive balance for a certain denomination or a positive number of LP shares. Transactions that would completely drain a pool are aborted.

Exiting with a swap requires paying both exit and swap fees.

#### Existing Exit Types

- `ExitPool`: The exit pool can be done in two ways 
  1. For a non-oracle pool or an oracle pool with an empty exit token denomination, it's not possible to exit entirely in a specific denomination. When a specified amount of shares needs to be exited, the shareRatio is calculated as follows: 
           ```shareRatio = exitShares / poolTotalShares``` The exit amount for each pool asset is then determined by multiplying this shareRatio by the respective amount of each asset in the pool.
  2. For an oracle pool with given exit token denom, exit token amount is calculated by calculating value of each LP token in the pool and then multiplied by exit share amount. 
     A `weightBreakingFee` is also applied which is calculated if distance from desired weight ratio increases after exiting the pool. `weightBreakingFee = WeightBreakingFeeMultiplier x [ ((weightIn x targetWeightOut) / (weightOut x targetWeightIn)) ^  WeightBreakingFeeExponent ]`
     where `weightIn = 1 - weightOut`, `weightOut = weight of asset in the pool in terms of value`, `targetWeightIn = 1 - targetWeightOut` and `targetWeightOut = weight of asset in the pool in terms of amount`

In both the methods an `exitFee` is applied before finally transferring it to the user. 
### Swap

When swapping assets within the pool, the input token is referred to as `tokenIn` and the output token as `tokenOut`. The module uses the following formula to calculate the number of tokens exchanged:

```
tokenBalanceOut * [1 - { tokenBalanceIn / (tokenBalanceIn + (1 - swapFee) * tokenAmountIn)} ^ (tokenWeightIn / tokenWeightOut)]
```

To reverse the calculation (i.e., given `tokenOut`), the following formula is used:

```
tokenBalanceIn * [{tokenBalanceOut / (tokenBalanceOut - tokenAmountOut)} ^ (tokenWeightOut / tokenWeightIn) - 1] / tokenAmountIn
```

#### Existing Swap Types

- `SwapExactAmountIn`
- `SwapExactAmountOut`

### Spot Price

The spot price, inclusive of the swap fee, is calculated using the formula:

```
spotPrice / (1 - swapFee)
```

Where `spotPrice` is defined as:

```
(tokenBalanceIn / tokenWeightIn) / (tokenBalanceOut / tokenWeightOut)
```

### Multi-Hop

Multi-hop logic is supported by the AMM module, allowing users to swap between multiple pools in a single transaction. The module calculates the optimal path for the swap, taking into account the swap fee and the spot price of each pool.

## Weights

Weights determine the distribution of assets within a pool. They are often expressed as ratios. For example, a 1:1 pool between "ETH" and "BTC" has a spot price of `#ETH in pool / #BTC in pool`. A 2:1 pool has a spot price of `2 * (#ETH in pool) / #BTC in pool`, which means fewer ETH reserves are required to maintain the same spot price, though it increases slippage.

Internally, weights are represented as numbers, with ratios computed as needed. Pools can be defined with weights up to `2^20`, and weight changes are managed with 30 bits of precision for smooth transitions.

## Network Parameters

Pools have the following configurable parameters:

| Key                      | Type                       |
| ------------------------ | -------------------------- |
| SwapFee                  | sdk.Dec                    |
| ExitFee                  | sdk.Dec                    |
| Weights                  | \*Weights                  |
| SmoothWeightChangeParams | \*SmoothWeightChangeParams |
| PoolCreationFee          | sdk.Coins                  |

### Parameter Definitions

- **SwapFee**: A percentage cut of all swaps that goes to the liquidity providers (LPs) of a pool.
- **ExitFee**: A fee applied when LPs remove their liquidity from the pool.
- **Weights**: Defines the asset weight ratios within the pool.
- **SmoothWeightChangeParams**: Allows for gradual weight adjustments over time.
- **PoolCreationFee**: The fee required to create a new pool, currently set to 20 USDC.
