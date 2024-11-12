#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Parameter module's query"
echo "Querying Params..."

$BINARY q parameter params --node $NODE

