#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running burner module's query"

echo "Querying params ..."
$BINARY q burner params --node $NODE

echo "Querying list-history ..."
$BINARY q burner list-history --node $NODE

echo "Querying show-history ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter timestamp and denom:"
    read timestamp denom
    $BINARY q burner show-history $timestamp $denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done