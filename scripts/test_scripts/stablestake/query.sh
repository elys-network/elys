#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Stablestake module's query"

echo "Querying Params ..."
$BINARY q stablestake params --node $NODE

echo "Querying borrow-ratio..."
$BINARY q stablestake borrow-ratio --node $NODE