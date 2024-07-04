<!--
order: 6
-->

# Slippage Model

In the Oracle Pool, to prevent bots from exploiting slippage by splitting trades into smaller pieces, we introduce the **stacked slippage model**. This model accumulates slippage within a block, making it less advantageous for bots to break down trades.

## Dynamic Weight Model

The Oracle Pool employs a dynamic weight model. The weight is updated at the start of the swap transactions execution and remains constant until the end of the block. As a result, slippage is recalculated based on the liquidity changes within the block. Multiple trades within a single block will cumulatively increase slippage, particularly for swaps in the same direction, leading to the "stacked slippage" model.

## Slippage Calculation

The slippage amount is calculated as follows:

$
SlippageAmount = Max(0, OracleOutAmount - BalancerOutAmount)
$

When $BalancerOutAmount$ exceeds $OracleOutAmount$, slippage is considered zero, not negative.

For trades with a negative slippage tolerance due to numerous opposite side trades, the value is set to 0%, aligning with the oracle price.

## Dynamic Weight Adjustment Criteria

To illustrate the dynamic weight adjustment:

- Suppose the pool contains 1000 ATOM and 10000 USDC, with the oracle providing a price of $11 for 1 ATOM.
- The pool's ratio is updated to 11:10 based on the USD value of each asset.

## Bot Prevention Mechanism

The stacked slippage model discourages bots by increasing slippage with consecutive same-direction swaps within a block:

- Initial swaps have minimal slippage based on the oracle price.
- Subsequent swaps face increased slippage, preventing lower overall slippage even if the trade is split.

For example, a bot attempting to swap 1000 ATOM for USDC with 1% slippage:

- Splitting the trade into 1000 swaps of 1 ATOM each would result in increasing slippage for each swap, cumulatively ensuring the overall slippage does not drop below 1%.

This mechanism ensures that splitting trades does not offer a better deal than a single trade within the block.

## Role of the Oracle

The oracle provides price and external liquidity data without additional information. This ensures price accuracy and consistency across trades.

## Dynamic Weight Model and Slippage

Oracle-based pools use the balancer slippage model to prevent arbitrage opportunities with other DEXs. The slippage is calculated from the balancer model, considering the weight and liquidity of the pool. When the weight is updated based on the oracle price, the pool's weight aligns with the oracle price, ensuring price consistency.

## Computational Considerations

Weights are updated when the first swap transaction is executed within a block, with all swap transactions executed at the end of the block. This does not directly impact users' swap gas fees since the transaction is only scheduled at this stage. The time complexity for executing swap transactions at the end of the block is O(n), where n is the number of swaps within the block.

## Handling Volatile Markets

In volatile markets, arbitrage opportunities arise due to price differences across platforms. The oracle provides the average price from these platforms and applies slippage on top, minimizing arbitrage opportunities with other platforms during volatile conditions.

## User Information on Slippage

Slippage is determined by liquidity, external liquidity, and oracle price. Users can estimate slippage using the balancer formula and set a minimum acceptable out amount, accounting for stacked slippage within the block. Similar to Uniswap, users can estimate the output amount, though it may vary due to other trades.
