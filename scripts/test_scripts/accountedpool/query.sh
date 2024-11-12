#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running accountedPool module's query"

echo "Querying list-accounted-pool ..."
$BINARY q poolaccounted list-accounted-pool --node $NODE

echo "Querying show-accounted-pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter index:"
    read index
    $BINARY q poolaccounted show-accounted-pool $index --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done