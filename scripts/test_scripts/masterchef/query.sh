#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running masterchef module's query"

echo "Querying Params ..."
$BINARY q masterchef params --node $NODE

echo "Querying Parlist-fee-infoams ..."
$BINARY q masterchef list-fee-info --node $NODE

echo "Querying APR ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the withdraw-type (from 0-5):"
    read earn_type
    $BINARY q masterchef apr $earn_type --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying external-incentive ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the reward-denom:"
    read reward_denom
    $BINARY q masterchef external-incentive $reward_denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying pool-aprs ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the pool-id:"
    read pool_id
    $BINARY q masterchef pool-aprs $pool_id --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying pool-info ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the pool-id:"
    read pool_id
    $BINARY q masterchef pool-info $pool_id --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying pool-reward-info ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the pool-id and reward-denom:"
    read pool_id reward_denom
    $BINARY q masterchef pool-reward-info $pool_id $reward_denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying pool-rewards ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the comma separated pool-id (i.e. 1,2):"
    read pool_ids
    $BINARY q masterchef pool-rewards $pool_ids --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying show-fee-info 2024-11-07 ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter date (i.e. 2024-11-07):"
    read date
    $BINARY q masterchef show-fee-info $date --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying stable-stake-apr ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter denom:"
    read denom
    $BINARY q masterchef stable-stake-apr $denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying user-pending-reward ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter address:"
    read address
    $BINARY q masterchef user-pending-reward $address --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying user-reward-info ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter address, pool-id and reward-denom:"
    read address pool_id reward_denom
    $BINARY q masterchef user-reward-info $address $pool_id $reward_denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done
