<!--
order: 8
-->

# Swap Transactions Reordering

This documentation describes the mechanism used to reorder swap transactions to minimize slippage and prevent MEV (Miner Extractable Value) attacks.

## Overview

To achieve the lowest possible slippage within a block, this module reorders swap requests into a specific sequence. This approach mitigates the risk of MEV attacks, which exploit transaction ordering. Instead of executing transactions as they are received from the validator, the requests are stored in a transient store during the endblocker and then executed in an optimal order.

## Key Considerations

When designing the reordering algorithm, the following factors are taken into account:

- Transactions that swap an exact amount in
- Transactions that swap for an exact amount out
- Multihop transactions
- Pools involving two or more tokens
- Algorithm complexity

## Algorithm

1. **Grouping Transactions**: Transactions are grouped by pool and token pairs (In -> Out).
   - For multihop transactions, only the first pool and the first in/out tokens are considered to reduce complexity.
2. **Selection and Execution**:

   - Select a random swap request.
   - Attempt to execute it in a cache context and check the resulting stacked slippage.
   - Check for an existing request in the opposite direction (same pool ID with opposite in/out tokens).
   - If such an opposite request exists, attempt to execute it in a cache context and check the resulting stacked slippage.
   - Apply the swap request that results in lower stacked slippage.
   - If either swap fails, discard the swap request without applying any changes.

3. **Iteration**:
   - Repeat the selection and execution process until all swap requests are processed.

### Algorithm Complexity

The complexity of the algorithm is `O(n)`, where `n` is the number of swap requests within the block.

### Algorithm Accuracy

This algorithm provides a semi-optimized solution, balancing accuracy and complexity.
