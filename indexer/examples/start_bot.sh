#!/bin/bash

# CLOB Matching Bot Startup Script
# This script loads configuration from environment or .env file

set -e

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check required environment variables
if [ -z "$BOT_MNEMONIC" ]; then
    echo "ERROR: BOT_MNEMONIC environment variable is required"
    echo "Please set it in your environment or create a .env file"
    exit 1
fi

# Set default values if not provided
export BOT_ENGINE_ID=${BOT_ENGINE_ID:-"matching-bot-$(hostname)-$$"}
export BOT_MATCHING_INTERVAL=${BOT_MATCHING_INTERVAL:-"500ms"}
export BOT_MAX_BATCH_SIZE=${BOT_MAX_BATCH_SIZE:-10}
export BOT_MIN_PROFIT_BPS=${BOT_MIN_PROFIT_BPS:-10}

export CHAIN_ID=${CHAIN_ID:-"elys-1"}
export CHAIN_RPC_ENDPOINT=${CHAIN_RPC_ENDPOINT:-"http://localhost:26657"}
export CHAIN_GRPC_ENDPOINT=${CHAIN_GRPC_ENDPOINT:-"localhost:9090"}
export CHAIN_GAS_PRICE=${CHAIN_GAS_PRICE:-"0.01uelys"}
export CHAIN_GAS_LIMIT=${CHAIN_GAS_LIMIT:-500000}
export CHAIN_DENOM=${CHAIN_DENOM:-"uelys"}

export INDEXER_URL=${INDEXER_URL:-"http://localhost:8080"}

# Default markets if not specified
export MARKETS=${MARKETS:-"[1,2,3]"}

# Expand environment variables in config file
CONFIG_FILE=${1:-bot_config_secure.yaml}
EXPANDED_CONFIG="/tmp/bot_config_expanded_$$.yaml"

echo "Expanding configuration file..."
envsubst < "$CONFIG_FILE" > "$EXPANDED_CONFIG"

# Build the bot if needed
if [ ! -f "./clob_matching_bot_v2" ]; then
    echo "Building bot..."
    go build -o clob_matching_bot_v2 clob_matching_bot_v2.go
fi

echo "Starting CLOB Matching Bot..."
echo "Engine ID: $BOT_ENGINE_ID"
echo "Chain ID: $CHAIN_ID"
echo "Indexer URL: $INDEXER_URL"
echo "Markets: $MARKETS"
echo "Min Profit: $BOT_MIN_PROFIT_BPS bps"

# Run the bot
./clob_matching_bot_v2 "$EXPANDED_CONFIG"

# Cleanup
rm -f "$EXPANDED_CONFIG"