# CLOB Redis Data Structure Design

## Overview
Redis will be used for high-performance caching and real-time order book management for the CLOB module. This design supports efficient order matching for external bots.

## Key Patterns

### 1. Order Book Structure
```
# Buy orders (sorted by price descending, then by counter ascending)
clob:orderbook:{market_id}:bids -> Sorted Set
  Score: price * 1e18 - counter
  Member: {order_id}

# Sell orders (sorted by price ascending, then by counter ascending)  
clob:orderbook:{market_id}:asks -> Sorted Set
  Score: price * 1e18 + counter
  Member: {order_id}

# Order details
clob:order:{order_id} -> Hash
  Fields:
    - market_id
    - owner
    - sub_account_id
    - order_type
    - price
    - amount
    - filled_amount
    - remaining_amount
    - status
    - created_at
    - block_height
```

### 2. Active Orders by Account
```
# User's active orders
clob:user:{owner}:orders -> Set
  Members: {order_id}

# User's orders by market
clob:user:{owner}:market:{market_id}:orders -> Set
  Members: {order_id}
```

### 3. Positions
```
# Open positions by market
clob:positions:market:{market_id} -> Set
  Members: {position_id}

# User's open positions
clob:user:{owner}:positions -> Set
  Members: {position_id}

# Position details
clob:position:{position_id} -> Hash
  Fields:
    - market_id
    - owner
    - sub_account_id
    - side
    - size
    - entry_price
    - mark_price
    - liquidation_price
    - margin
    - unrealized_pnl
    - updated_at
```

### 4. Market Data
```
# Market information
clob:market:{market_id} -> Hash
  Fields:
    - ticker
    - base_asset
    - quote_asset
    - last_price
    - best_bid
    - best_ask
    - volume_24h
    - trades_24h
    - open_interest
    - funding_rate
    - next_funding_time

# Price tickers (for quick access)
clob:tickers -> Hash
  Field: {market_id}
  Value: JSON with price data

# Recent trades
clob:trades:market:{market_id} -> List (capped at 1000)
  Elements: JSON trade data

# Trade history by user
clob:user:{owner}:trades -> List (capped at 1000)
  Elements: JSON trade data
```

### 5. Order Matching Support
```
# Best bid/ask prices for quick matching
clob:market:{market_id}:best_bid -> String (price)
clob:market:{market_id}:best_ask -> String (price)

# Market order queue
clob:market:{market_id}:market_orders -> List
  Elements: {order_id}

# Order locks (for atomic matching)
clob:order:lock:{order_id} -> String (with TTL)
  Value: matching_engine_id
```

### 6. Real-time Updates
```
# WebSocket subscriptions
clob:ws:subscriptions -> Hash
  Field: {client_id}
  Value: JSON subscription filters

# Update streams (using Redis Streams)
clob:stream:orders -> Stream
clob:stream:trades -> Stream
clob:stream:positions -> Stream
clob:stream:orderbook -> Stream
```

### 7. Statistics and Analytics
```
# 24h high/low
clob:market:{market_id}:high24h -> String
clob:market:{market_id}:low24h -> String

# Funding rate history
clob:funding:{market_id} -> Sorted Set
  Score: timestamp
  Member: JSON with rate data

# Liquidation queue
clob:liquidations:queue -> Sorted Set
  Score: liquidation_price
  Member: {position_id}
```

## TTL Strategy

- Order locks: 5 seconds
- Completed orders: 1 hour
- Trade history: 24 hours
- Price tickers: No TTL (always current)
- Order book entries: No TTL (removed on fill/cancel)

## Atomic Operations

For order matching, use Redis transactions (MULTI/EXEC) or Lua scripts:

```lua
-- Example: Atomic order match
local buy_order = redis.call('HGETALL', KEYS[1])
local sell_order = redis.call('HGETALL', KEYS[2])

if buy_order.price >= sell_order.price then
  -- Execute trade
  redis.call('HINCRBY', KEYS[1], 'filled_amount', trade_amount)
  redis.call('HINCRBY', KEYS[2], 'filled_amount', trade_amount)
  -- Update order book
  -- Publish trade event
end
```

## Performance Considerations

1. Use pipelining for bulk operations
2. Implement connection pooling
3. Use Redis Cluster for horizontal scaling
4. Monitor memory usage and implement eviction policies
5. Use Redis Streams for event streaming instead of pub/sub for persistence