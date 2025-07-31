# CLOB Integration Guide

## Overview

The CLOB (Central Limit Order Book) integration provides comprehensive indexing and real-time data access for the Elys CLOB module. This includes order book management, trade history, position tracking, and support for external order matching bots.

## Architecture

### Components

1. **PostgreSQL Tables**
   - `clob_markets` - Market configuration and parameters
   - `clob_orders` - All orders (active and historical)
   - `clob_trades` - Executed trades
   - `clob_positions` - Perpetual positions
   - `clob_order_book_snapshots` - Point-in-time order book states
   - `clob_funding_rates` - Funding rate history
   - `clob_liquidations` - Liquidation events
   - `clob_market_stats` - Aggregated market statistics

2. **Redis Cache**
   - Real-time order book management
   - Active orders and positions
   - Market data and tickers
   - Order locking for matching
   - WebSocket pub/sub for real-time updates

3. **Event Processing**
   - CLOB event parser for blockchain events
   - Real-time order book updates
   - Position tracking
   - Trade execution handling

## API Endpoints

### Market Data

```
GET /api/v1/clob/markets
GET /api/v1/clob/markets/{marketId}
GET /api/v1/clob/markets/{marketId}/stats
```

### Order Book

```
GET /api/v1/clob/markets/{marketId}/orderbook
GET /api/v1/clob/markets/{marketId}/orderbook/snapshot
GET /api/v1/clob/markets/{marketId}/best
```

### Orders

```
GET /api/v1/clob/orders/active
GET /api/v1/clob/orders/{orderId}
GET /api/v1/clob/users/{address}/orders
```

### Trades

```
GET /api/v1/clob/markets/{marketId}/trades
GET /api/v1/clob/users/{address}/trades
```

### Positions

```
GET /api/v1/clob/users/{address}/positions
GET /api/v1/clob/positions/{positionId}
```

### Order Matching Support

```
GET /api/v1/clob/matching/orders
POST /api/v1/clob/matching/lock
POST /api/v1/clob/matching/unlock
```

## Order Matching for External Bots

The indexer provides APIs specifically designed for external matching engines:

### 1. Get Orders for Matching

```bash
GET /api/v1/clob/matching/orders?market_id=1
```

Returns active buy and sell orders grouped by side:

```json
{
  "market_id": 1,
  "timestamp": "2024-01-01T00:00:00Z",
  "buy_orders": [...],
  "sell_orders": [...]
}
```

### 2. Lock Orders

Before matching, lock orders to prevent double-spending:

```bash
POST /api/v1/clob/matching/lock
{
  "order_id": 12345,
  "matching_engine_id": "bot-1",
  "ttl_seconds": 5
}
```

### 3. Unlock Orders

After processing (success or failure), unlock orders:

```bash
POST /api/v1/clob/matching/unlock
{
  "order_id": 12345,
  "matching_engine_id": "bot-1"
}
```

### Example Matching Bot

See `/examples/clob_matching_bot.go` for a complete example of:
- Connecting to the indexer
- Fetching order books
- Finding matching opportunities
- Locking orders
- Processing matches

## WebSocket Subscriptions

Real-time updates are available via WebSocket:

### Order Updates
```
Channel: clob:updates:orders:{owner_address}
```

### Trade Updates
```
Channel: clob:updates:trades:{market_id}
```

### Position Updates
```
Channel: clob:updates:positions:{owner_address}
```

### Order Book Updates
```
Channel: clob:updates:orderbook:{market_id}
```

## Redis Data Structure

### Order Book
```
clob:orderbook:{market_id}:bids - Sorted set of buy orders
clob:orderbook:{market_id}:asks - Sorted set of sell orders
```

### Order Details
```
clob:order:{order_id} - Hash of order details
```

### Best Prices
```
clob:market:{market_id}:best_bid
clob:market:{market_id}:best_ask
```

### Order Locks
```
clob:order:lock:{order_id} - Lock for order matching
```

## Performance Considerations

1. **Order Book Aggregation**
   - Runs every 5 seconds by default
   - Configurable via `OrderBookUpdateInterval`
   - Maintains top 20 levels by default

2. **Caching Strategy**
   - Active orders cached in Redis
   - Best bid/ask updated on every order change
   - TTL on completed orders (1 hour)

3. **Scaling**
   - Redis Cluster for horizontal scaling
   - Multiple indexer instances with shared Redis
   - Connection pooling for database queries

## Configuration

Add to `config.yaml`:

```yaml
indexer:
  order_book_update_interval: 5s
  event_buffer_size: 1000
  worker_count: 10

redis:
  address: "localhost:6379"
  password: ""
  db: 0
  pool_size: 100
```

## Testing

1. Run the indexer with CLOB support
2. Create test orders via blockchain
3. Monitor order book updates
4. Test matching bot with small amounts
5. Verify trades are recorded correctly

## Security Considerations

1. **Order Locking**
   - TTL prevents permanent locks
   - Engine ID prevents unauthorized unlocks
   - Atomic operations via Redis

2. **Rate Limiting**
   - Implement rate limits on API endpoints
   - Monitor for suspicious matching patterns

3. **Data Validation**
   - Verify order ownership before matching
   - Check order status before execution
   - Validate price/amount calculations