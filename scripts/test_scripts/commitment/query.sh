#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Commitment module's query"

echo "Querying Params ..."
$BINARY q commitment params --node $NODE

echo "Querying ShowCommitment..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the creator:"
    read creator
    $BINARY q commitment show-commitments $creator --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying ShowCommittedTokensLocked..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q commitment committed-tokens-locked $address --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying number-of-commitments ..."
$BINARY q commitment number-of-commitments --node $NODE

echo "Querying commitment-vesting-info..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q commitment commitment-vesting-info $address --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done