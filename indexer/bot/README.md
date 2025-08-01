# CLOB Matching Bots

This directory contains two implementations of CLOB (Central Limit Order Book) matching bots, each in its own directory to avoid naming conflicts:

1. **Redis Bot** (`redis_bot/`) - Uses Redis for order storage and matching
2. **PostgreSQL Bot** (`postgres_bot/`) - Uses PostgreSQL for order storage and matching

## Directory Structure

```
bot/
├── redis_bot/
│   └── main.go      # Redis bot implementation
├── postgres_bot/
│   └── main.go      # PostgreSQL bot implementation
├── bot_config.yaml  # Shared configuration
├── run_bots.sh      # Script to run bots
└── README.md        # This file
```

## Features

Both bots:
- Monitor specified markets for matching opportunities
- Use the same mnemonic for wallet derivation
- Match buy and sell orders when prices cross
- Implement order locking to prevent double-matching
- Run continuously with configurable intervals
- Include test data for immediate demonstration

## Prerequisites

Make sure the indexer services are running:
```bash
cd ../
docker-compose up -d
```

This starts:
- PostgreSQL on port 5433
- Redis on port 6380
- Indexer service on port 8080

## Running the Bots

### Run Redis Bot Only
```bash
./run_bots.sh redis
```

### Run PostgreSQL Bot Only
```bash
./run_bots.sh postgres
```

### Run Both Bots
```bash
./run_bots.sh both
```

### Run Directly
```bash
# Redis bot
cd redis_bot && go run main.go

# PostgreSQL bot
cd postgres_bot && go run main.go
```

## Configuration

The bots use:
- Redis URL: `redis://localhost:6380`
- PostgreSQL URL: `postgres://indexer:123123123@localhost:5433/elys_indexer?sslmode=disable`
- Mnemonic: Configured in the code
- Market ID: 1 (by default)

## How They Work

### Redis Bot
- Stores orders in Redis sorted sets (by price)
- Uses Redis SETNX for distributed locking
- Stores matches in Redis with 24-hour TTL
- Seeds test data on startup

### PostgreSQL Bot
- Queries orders from the `clob_orders` table
- Uses database transactions for consistency
- Stores matches in the `clob_matches` table
- Creates tables automatically if they don't exist
- Seeds test data on startup

## Database Schema

The PostgreSQL bot creates these tables:
- `clob_orders` - Active orders
- `clob_matches` - Executed matches

## Development

To modify the bots:
1. Edit the respective `main.go` files
2. Run with `go run main.go` from the bot directory

## Notes

- Both bots operate independently and can run simultaneously
- They use the same mnemonic but different storage backends
- Each bot is in its own directory to avoid naming conflicts
- Uses the parent indexer module's dependencies
- In production, implement proper transaction signing and blockchain submission
- Add proper error handling and monitoring for production use