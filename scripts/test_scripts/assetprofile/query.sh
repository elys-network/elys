#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Assetprofile module's query"

echo "Querying Params ..."
$BINARY q assetprofile params --node $NODE

echo "Querying All Entry ..."
$BINARY q assetprofile list-entry --node $NODE

echo "Querying Entry item by base denom ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter base_denom:"
    read BASE_DENOM
    $BINARY q assetprofile show-entry $BASE_DENOM --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying Entry items by denom ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter denom:"
    read DENOM
    $BINARY q assetprofile show-entry-by-denom $DENOM --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done
