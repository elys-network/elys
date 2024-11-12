#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Clock module's query"

echo "Querying contracts ..."
$BINARY q clock contracts --node $NODE

echo "Querying params ..."
$BINARY q clock params --node $NODE