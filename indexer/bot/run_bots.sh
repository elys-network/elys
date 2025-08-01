#!/bin/bash

# Script to run both bots

echo "Starting CLOB Matching Bots..."

# Set environment variables
export REDIS_URL="redis://localhost:6380"
export DATABASE_URL="postgres://indexer:123123123@localhost:5433/elys_indexer?sslmode=disable"

# Check if services are running
echo "Checking Redis connection..."
redis-cli -p 6380 ping > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "Error: Redis is not running on port 6380"
    exit 1
fi

echo "Checking PostgreSQL connection..."
PGPASSWORD=123123123 psql -h localhost -p 5433 -U indexer -d elys_indexer -c "SELECT 1" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "Error: PostgreSQL is not running on port 5433"
    exit 1
fi

# Run bots based on argument
case "$1" in
    redis)
        echo "Starting Redis bot..."
        cd redis_bot && go run main.go
        ;;
    postgres)
        echo "Starting PostgreSQL bot..."
        cd postgres_bot && go run main.go
        ;;
    both)
        echo "Starting both bots..."
        (cd redis_bot && go run main.go) &
        REDIS_PID=$!
        (cd postgres_bot && go run main.go) &
        POSTGRES_PID=$!
        
        echo "Redis bot PID: $REDIS_PID"
        echo "PostgreSQL bot PID: $POSTGRES_PID"
        
        # Wait for both processes
        wait $REDIS_PID $POSTGRES_PID
        ;;
    *)
        echo "Usage: $0 {redis|postgres|both}"
        exit 1
        ;;
esac