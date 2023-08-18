# Swap transactions reordering

To provide lowest slippage as possible on the block, module sort the swap requests into a specific order.

This removes possibility of MEV attack on swap transactions.

Instead of executing messages directly, on the endblocker, we store the requests on transient store and execute in specific order to minimize maximum slippage.

## To consider on algorithm

Swap exact amount in txs
Swap exact amount out txs
Multihop txs
2+ token pools
Algorithm complexity

## Algorithm

- Group transactions by pool and token (In -> Out) pairs
  - Mulithop txs just consider first pool and first in/out tokens to avoid complexity
- Select a random swap request
  - Try execution on cache context, and check stacked slippage
  - Check if opposite direction request exists (Same pool id with opposite in/out tokens)
  - If opposite direction request exists, try execution on cache context, and check stacked slippage
  - Apply the swap request which as lower stacked slippage
  - If one of the swaps fail, not apply any changes and remove the swap request
- Repeat the process until the swap requests run-out

### Algorithm Complexity

`O(n)` - where `n` is number of swap requests on the block

### Algorithm Accuracy

This is semi-optimized solution while keeping the algorithm complexity low.
