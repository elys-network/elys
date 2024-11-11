#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running estaking module's query"

echo "Querying Params ..."
$BINARY q estaking params --node $NODE

echo "Querying invariant ..."
$BINARY q estaking invariant --node $NODE

echo "Querying rewards ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address]:"
    read address
    $BINARY q estaking rewards $address
    echo "Want to make another query, true/false"
    read ASK_USER
done


