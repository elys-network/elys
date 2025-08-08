# Memory OrderBook Optimization

## Overview
The in-memory orderbook has been optimized to avoid loading the entire orderbook from state on every block. Instead, it maintains a persistent synchronized copy in memory with efficient price-level indexing.

## Key Improvements

### 1. Price-Level Indexing
- Orders are grouped by price tick (multiple of 1,000,000) 
- Buy orders: `marketId -> priceTick -> [orders]` (sorted descending by price)
- Sell orders: `marketId -> priceTick -> [orders]` (sorted ascending by price)
- Separate sorted price level arrays for efficient iteration

### 2. Persistent Memory State
- Singleton `MemoryOrderBook` instance in the Keeper
- Initialized once on startup from blockchain state
- Stays synchronized with all order operations

### 3. Automatic Synchronization
All order operations automatically update the memory orderbook:
- `SetOrder()`: Adds new orders or updates existing ones
- `DeleteOrder()`: Removes orders from memory
- `BatchDeleteOrders()`: Batch removal from memory
- Order execution automatically updates filled amounts

### 4. Efficient Matching Algorithm
```go
// Pseudocode for matching
for each buy_price_level (descending):
    for each sell_price_level (ascending):
        if buy_price >= sell_price:
            match orders at these levels
        else:
            break (no more crosses possible)
```

## Performance Benefits

### Before (Loading from State)
- **Every block**: O(N) to load all orders from state
- **Memory**: Creates new data structures each time
- **I/O**: Heavy disk reads on every vote extension

### After (Persistent Memory)
- **Startup**: O(N) one-time load
- **Per operation**: O(log P) where P = number of price levels
- **Memory**: Reuses existing structures
- **I/O**: No disk reads during matching

## Data Structures

```go
type MemoryOrderBook struct {
    // Price-level maps for efficient access
    buyOrdersByPrice  map[marketId]map[priceTick]*OrderList
    sellOrdersByPrice map[marketId]map[priceTick]*OrderList
    
    // Sorted price levels for iteration
    buyPriceLevels  map[marketId][]priceTick  // descending
    sellPriceLevels map[marketId][]priceTick  // ascending
    
    // Quick lookup by order ID
    ordersByID map[marketId]map[orderId]*Order
}
```

## Edge Cases Handled

1. **Full Orderbook Scan**: In worst case where all orders need checking, the algorithm efficiently iterates through price levels rather than individual orders.

2. **Empty Price Levels**: Automatically removed when last order at a price is deleted.

3. **Order Updates**: Filled amounts are tracked, and fully filled orders are automatically removed.

4. **Market Creation**: New markets are automatically initialized in memory when orders are placed.

## Usage in Vote Extensions

```go
func (k Keeper) GetOperationsToPropose(ctx, marketIds, maxMatchesPerMarket) {
    // No loading from state needed!
    // Memory orderbook is already synchronized
    
    for each marketId:
        matched = memoryOrderBook.MatchOrders(marketId, maxMatches)
        // Convert to vote extension format
}
```

## Synchronization Points

1. **Keeper Creation**: Empty memory orderbook created
2. **First Vote Extension**: Lazy initialization from state
3. **Order Placement**: Added to memory
4. **Order Execution**: Updated in memory
5. **Order Deletion**: Removed from memory

## Testing Recommendations

1. **Correctness**: Verify memory state matches blockchain state
2. **Performance**: Benchmark matching with large orderbooks
3. **Concurrency**: Test with concurrent order operations
4. **Recovery**: Test state recovery after restart

## Future Optimizations

1. **Incremental Updates**: Track only changed orders between blocks
2. **Parallel Matching**: Match multiple markets concurrently
3. **Memory Pooling**: Reuse order objects to reduce GC pressure
4. **Price Level Caching**: Cache best bid/ask for quick access