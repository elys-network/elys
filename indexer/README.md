# Elys Chain Indexer

A PostgreSQL and Redis-based indexer for the Elys blockchain, specifically designed to index TradeSheild and Perpetual module events. The indexer provides real-time WebSocket updates for order books, trades, and positions.

## Features

- **Event Indexing**: Indexes events from TradeSheild and Perpetual modules
- **Order Book Management**: Maintains real-time order book data for spot and perpetual markets
- **Trade History**: Records all executed trades with full details
- **Position Tracking**: Tracks perpetual positions (MTPs) with real-time updates
- **WebSocket API**: Provides real-time updates for order books, trades, and positions
- **Redis Caching**: High-performance caching layer for frequently accessed data
- **PostgreSQL Storage**: Persistent storage for all indexed data

## Architecture

```
┌─────────────────┐
│   Elys Chain    │
│  (CometBFT RPC) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐      ┌─────────────────┐
│     Indexer     │◄────►│   PostgreSQL    │
│  (Event Parser) │      │   (Storage)     │
└────────┬────────┘      └─────────────────┘
         │
         ▼
┌─────────────────┐      ┌─────────────────┐
│  Redis Cache    │◄────►│WebSocket Server │
│  (Pub/Sub)      │      │  (Real-time)    │
└─────────────────┘      └─────────────────┘
```

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (optional)

## Installation

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/elys-network/elys
cd elys/indexer
```

2. Configure the indexer:
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your chain endpoints
```

3. Start the services:
```bash
make docker-up
```

### Manual Installation

1. Install dependencies:
```bash
go mod download
```

2. Set up PostgreSQL database:
```bash
psql -U postgres -c "CREATE DATABASE elys_indexer;"
psql -U postgres -d elys_indexer -f sql/schema.sql
```

3. Configure Redis (default configuration should work)

4. Build and run:
```bash
make build
./elys-indexer -c config.yaml
```

## Configuration

The indexer is configured via `config.yaml`:

```yaml
chain:
  rpc_endpoint: "http://localhost:26657"
  grpc_endpoint: "localhost:9090"
  chain_id: "elys-1"
  start_height: 0  # 0 means start from latest

database:
  host: "localhost"
  port: 5432
  user: "indexer"
  password: "indexer_password"
  database: "elys_indexer"

redis:
  addr: "localhost:6379"

websocket:
  listen_addr: ":8080"
```

## WebSocket API

### Connection

Connect to the WebSocket server at `ws://localhost:8080/ws`

### Message Format

All messages use JSON format:

```json
{
  "type": "subscribe|unsubscribe|ping|pong|update|error",
  "channel": "channel_name",
  "data": {}
}
```

### Subscription Types

1. **Order Book Updates**
```json
{
  "type": "subscribe",
  "data": {
    "type": "order_book",
    "filters": {
      "asset_pair": "ATOM/USDC"
    }
  }
}
```

2. **Trade Updates**
```json
{
  "type": "subscribe",
  "data": {
    "type": "trades",
    "filters": {
      "asset": "ATOM"
    }
  }
}
```

3. **Order Updates**
```json
{
  "type": "subscribe",
  "data": {
    "type": "orders",
    "filters": {
      "owner": "elys1..."
    }
  }
}
```

4. **Position Updates**
```json
{
  "type": "subscribe",
  "data": {
    "type": "positions",
    "filters": {
      "owner": "elys1..."
    }
  }
}
```

## Development

### Running Tests
```bash
make test
```

### Running in Development Mode
```bash
make dev
```

### Code Formatting
```bash
make fmt
```

### Linting
```bash
make lint
```

## API Endpoints

### HTTP Endpoints

- `GET /health` - Health check endpoint

### WebSocket Endpoints

- `ws://localhost:8080/ws` - WebSocket connection endpoint

## Database Schema

The indexer uses the following main tables:

- `spot_orders` - TradeSheild spot orders
- `perpetual_orders` - TradeSheild perpetual orders
- `perpetual_positions` - Perpetual module positions (MTPs)
- `trades` - Executed trades from both modules
- `order_book_snapshots` - Periodic order book snapshots

## Monitoring

The indexer exposes Prometheus metrics on port 9090:

- `indexer_blocks_processed_total` - Total blocks processed
- `indexer_events_processed_total` - Total events processed
- `indexer_processing_duration_seconds` - Processing duration histogram
- `indexer_websocket_connections` - Current WebSocket connections

## Troubleshooting

### Common Issues

1. **Connection refused to chain**
   - Verify RPC/gRPC endpoints are correct
   - Check if the chain node is running

2. **Database connection failed**
   - Ensure PostgreSQL is running
   - Check database credentials

3. **WebSocket connection drops**
   - Check client ping/pong implementation
   - Verify network stability

### Logs

View logs using:
```bash
# Docker
make docker-logs

# Local
./elys-indexer -c config.yaml -l debug
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.