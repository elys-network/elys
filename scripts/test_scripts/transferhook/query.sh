#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running transferhook module's query"

echo "Querying Params ..."
$BINARY q transferhook params --node $NODE


