# CLOB Matching Bot V2

A production-ready order matching bot for the Elys CLOB module that can execute batch transactions on-chain.

## Features

- **Batch Processing**: Matches multiple orders in a single transaction to save gas
- **Profit Optimization**: Only executes profitable matches based on configurable minimum profit threshold
- **Secure Configuration**: Supports environment variables for sensitive data
- **Order Locking**: Prevents race conditions with atomic order locking
- **Multi-Market Support**: Monitor and match orders across multiple markets simultaneously
- **Automatic Account Management**: Fetches account info and manages nonces automatically

## Prerequisites

1. Elys node with gRPC enabled
2. Indexer API running with CLOB support
3. A funded wallet for the bot
4. Go 1.21 or higher

## Setup

### 1. Generate a Bot Wallet

```bash
# Generate a new wallet
elysd keys add bot

# Or recover from mnemonic
elysd keys add bot --recover

# Get the address
elysd keys show bot -a

# Fund the wallet with ELYS tokens for gas
```

### 2. Configure the Bot

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# REQUIRED: Your bot wallet mnemonic
BOT_MNEMONIC="your mnemonic phrase here"

# Profit threshold (basis points, 100 = 1%)
BOT_MIN_PROFIT_BPS=10

# Markets to monitor
MARKETS=[1,2,3]

# Chain endpoints
CHAIN_RPC_ENDPOINT=http://localhost:26657
CHAIN_GRPC_ENDPOINT=localhost:9090

# Indexer endpoint
INDEXER_URL=http://localhost:8080
```

### 3. Build and Run

Using the startup script:

```bash
# Make the script executable
chmod +x start_bot.sh

# Run with environment variables
./start_bot.sh

# Or run with a specific config file
./start_bot.sh bot_config.yaml
```

Manual build and run:

```bash
# Build
go build -o clob_matching_bot_v2 clob_matching_bot_v2.go

# Run with config file
./clob_matching_bot_v2 bot_config.yaml
```

## Configuration Options

### Bot Settings

| Variable | Description | Default |
|----------|-------------|---------|
| `BOT_ENGINE_ID` | Unique identifier for this bot instance | `matching-bot-{hostname}-{pid}` |
| `BOT_MNEMONIC` | Wallet mnemonic (24 words) | Required |
| `BOT_MATCHING_INTERVAL` | How often to check for matches | `500ms` |
| `BOT_MAX_BATCH_SIZE` | Max matches per transaction | `10` |
| `BOT_MIN_PROFIT_BPS` | Minimum profit in basis points | `10` |

### Chain Settings

| Variable | Description | Default |
|----------|-------------|---------|
| `CHAIN_ID` | Chain ID | `elys-1` |
| `CHAIN_RPC_ENDPOINT` | Tendermint RPC endpoint | `http://localhost:26657` |
| `CHAIN_GRPC_ENDPOINT` | Cosmos gRPC endpoint | `localhost:9090` |
| `CHAIN_GAS_PRICE` | Gas price per unit | `0.01uelys` |
| `CHAIN_GAS_LIMIT` | Gas limit per transaction | `500000` |

## How It Works

1. **Order Discovery**: The bot queries the indexer API for active orders in configured markets
2. **Match Finding**: Identifies crossing orders where buy price ≥ sell price
3. **Profit Calculation**: Calculates profit margin and filters based on minimum threshold
4. **Order Locking**: Locks matched orders via the indexer to prevent double-spending
5. **Batch Creation**: Groups multiple matches into a single transaction
6. **Transaction Execution**: Signs and broadcasts the batch match transaction
7. **Unlock Orders**: Releases locks after transaction processing

## Profit Calculation

The bot calculates profit as:

```
Profit (bps) = (Buy Price - Sell Price) / Sell Price * 10000
```

For example:
- Buy order at 100 USDC
- Sell order at 99.9 USDC
- Profit = (100 - 99.9) / 99.9 * 10000 = 10 bps (0.1%)

## Transaction Format

The bot creates `MsgBatchMatch` transactions with the following structure:

```protobuf
message MsgBatchMatch {
  string creator = 1;
  repeated MatchOperation operations = 2;
}

message MatchOperation {
  OrderKey buy_order_key = 1;
  OrderKey sell_order_key = 2;
  string quantity = 3;
  string price = 4;
}
```

## Monitoring

The bot logs all activities:

```
2024/01/01 00:00:00 Starting matching bot v2 with address: elys1...
2024/01/01 00:00:01 Processing batch 1/1 with 3 matches
2024/01/01 00:00:01 Matched: Buy #123 @ 100.5 <-> Sell #456 @ 100.0, Amount: 10.5, Profit: 50 bps
2024/01/01 00:00:02 Batch match transaction submitted: 0x123...
```

## Security Considerations

1. **Mnemonic Storage**: Never commit mnemonics to version control. Use environment variables or secure key management.
2. **Network Security**: Run the bot in a secure network environment with proper firewall rules.
3. **Rate Limiting**: The indexer should implement rate limiting to prevent abuse.
4. **Monitoring**: Set up alerts for failed transactions and unusual activity.
5. **Gas Management**: Monitor gas prices and adjust configuration as needed.

## Troubleshooting

### Common Issues

1. **"Failed to lock order"**: Order is already being processed by another bot
2. **"Transaction failed: insufficient fees"**: Increase gas price or limit
3. **"Account sequence mismatch"**: Bot will auto-retry with correct sequence
4. **"Invalid mnemonic"**: Check that mnemonic is 24 words and properly formatted

### Debug Mode

Set log level for more details:

```bash
export LOG_LEVEL=debug
./start_bot.sh
```

## Production Deployment

For production use:

1. Use a dedicated server or container
2. Set up process management (systemd, supervisor)
3. Configure log rotation
4. Set up monitoring and alerting
5. Use secure key management (HSM, Vault)
6. Implement circuit breakers for API calls
7. Set up database for transaction history

## Example Systemd Service

```ini
[Unit]
Description=CLOB Matching Bot
After=network.target

[Service]
Type=simple
User=elys
WorkingDirectory=/opt/elys-bot
EnvironmentFile=/opt/elys-bot/.env
ExecStart=/opt/elys-bot/start_bot.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## Development

To modify the bot:

1. Update `clob_matching_bot_v2.go`
2. Add new configuration options to `Config` struct
3. Update `bot_config.yaml` and `.env.example`
4. Test thoroughly on testnet before mainnet

## License

Apache 2.0