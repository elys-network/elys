# Slippage model

On oracle pool, to prevent bots from splitting the trades into smaller pieces to reduce slippage applied on the trade, we introduce the stacked slippage model where slippage is stacked for a block.

Oracle pool has a dynamic weight model, the weight is updated before starting the swap transactions execution and does not change until the end of the block.

As pool weight is updated per block (not per swap tx), slippage is changed on the block based on liquidity change on the block.
Multiple trades within a single block will cumulatively increase the slippage. When same direction swaps are executed, slippage is being increased for later swaps. We call this model as `stacked slippage` model.

## Slippage calculation

```
SlippageAmount = Max(0, OracleOutAmount - BalancerOutAmount)
```

When `BalancerOutAmount` is bigger than `OracleOutAmount`, slippage is considered as 0 not negative.

For trades that has negative slippage tolerance for lots of opposite side trades, the value is set to 0% and it returns exactly same as oracle price.

## How exactly is the weight updated before the swap transactions execution? What is the criteria or the algorithm behind this dynamic weight adjustment?

Let's say pool has 1000 ATOM and 10000 USDC. And oracle has provided $11 for 1 ATOM data.
The ratio for the pool is updated to 11:10 based on USD value of each asset on the pool.

## How exactly does the stacked slippage model prevent bots

Since slippage is stacked, the first swap would have pretty small slippage based on oracle price. But as time goes, the slippage is increased if same direction swaps are executed.

E.g. if a bot would like to swap 1000ATOM for USDC where slippage is 1%.
If the bot split this trade into 1000 pieces of 1 ATOM -> USDC swap, first swap slippage is 0.001%. But it's cumulatively increased as next 1 ATOM -> USDC swaps are executed. And on the block, the user can't get lower than 1% slippage for overall 1000 ATOM.

The bot could split the swap and run on next block, but on next block, price should have been changed and therefore, bot can't estimate output amount exactly on current block.

And therefore, we can ensure that dust amount split attack is not generating better deal than 1 time swap on the block.

## How does the oracle play into this

Oracle is providing price data and external liquidity data and it's not providing any other information.

## How does the dynamic weight model relate to slippage

Oracle based pools use balancer slippage model to prevent arbitrage opportunity on swaps with other DEXs.
It is calculating slippage from balancer model - weight and liquidity on the pool.

When the weight is changed based on oracle price, pool weights are updated based on oracle price, and for that, price on the pool is exactly same as oracle price.

## Are there any gas or computation considerations when updating weights or calculating slippage, especially if these operations are done frequently?

Weight are updated when the first swap transaction is executed on the block. And all the swap transactions are executed on the endblocker.
This won't be directly affecting users' swap gas fees where transaction is only scheduled at the stage.

The overall time complexity on the endblocker for the swap transactions take `O(n)` where `n` is the number of swaps on the block.

## How does the model account for volatile markets or rapid price changes?

On volatile markets, there are a lot of arbitrage opportunities between platforms for price difference.
The oracle is providing average price of those platforms, and applying slippage on top of it as well.
Therefore, this model minimize the arbitrage opportunity with other platforms on volatile markets.

## How are users informed about the current slippage, especially if it's stacked and changing within a block?

Slippage is based on liquidity, external liquidity and oracle price. The slippage can be calculated from balancer formula. Users set minimum out amount from slippage estimation + additional acceptable stacked slippage for the block. Uniswap is also doing same thing, you can estimate the out amount but it's not an exact amount, it can be affected by other trades.
