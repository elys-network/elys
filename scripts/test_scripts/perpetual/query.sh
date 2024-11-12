#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running perpetual module's query"

echo "Querying Params ..."
$BINARY q perpetual params --node $NODE


echo "Querying Close estimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] [position-id] [closing-amount] [closing-price]:"
    read address position_id closing_amount closing_price
    $BINARY q perpetual close-estimation $address $position_id $closing_amount $closing_price
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying Open estimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [position] [leverage] [trading-asset] [collateral] [pool-id]:"
    read position leverage trading_asset collateral pool_id
    $BINARY q perpetual open-estimation $position $leverage $trading_asset $collateral $pool_id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-all-to-pay ..."
$BINARY q perpetual get-all-to-pay --node $NODE

echo "Querying get-positions ..."
$BINARY q perpetual get-positions --node $NODE

echo "Querying get-status ..."
$BINARY q perpetual get-status --node $NODE

echo "Querying get-whitelist ..."
$BINARY q perpetual get-whitelist --node $NODE

echo "Querying list-pool ..."
$BINARY q perpetual list-pool --node $NODE

echo "Querying is-whitelisted ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q perpetual is-whitelisted $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-mtp ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address and position_id:"
    read address position_id
    $BINARY q perpetual get-mtp $address $position_id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-positions-by-pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the amm_pool_id:"
    read amm_pool_id
    $BINARY q perpetual get-positions-by-pool $amm_pool_id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-positions-for-address ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q perpetual get-positions-for-address $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying show-pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the index:"
    read index
    $BINARY q perpetual show-pool $index
    echo "Want to make another query, true/false"
    read ASK_USER
done

