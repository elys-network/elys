#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running epochs module's query"

echo "Querying epoch-infos ..."
$BINARY q epochs epoch-infos --node $NODE

echo "Querying current-epoch ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [identifier]:"
    read identifier
    $BINARY q epochs current-epoch $identifier
    echo "Want to make another query, true/false"
    read ASK_USER
done


