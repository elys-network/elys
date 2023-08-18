# Slippage model

On oracle pool, to prevent bots from splitting the trades into smaller pieces to reduce slippage applied on the trade, we introduce the stacked slippage model where slippage is stacked for a block.

Oracle pool has a dynamic weight model, the weight is updated before starting the swap transactions execution and does not change until the end of the block.

For trades where slippage tolerance is a negative value for lots of opposite side trades, the slippage tolerance value is set to 0%.
