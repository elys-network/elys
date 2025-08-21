#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running leveragelp module's query"

echo "Querying Params ..."
$BINARY q leveragelp params --node $NODE

echo "Querying Close estimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [owner_address] [pool_id] [lp_amount]:"
    read owner pool_id lp_amount
    $BINARY q leveragelp close-estimation $owner $pool_id $lp_amount
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying committed-tokens-locked ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q leveragelp committed-tokens-locked $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-position ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address and id:"
    read address id
    $BINARY q leveragelp get-position $address $id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-positions-by-pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the amm pool id:"
    read id
    $BINARY q leveragelp get-positions-by-pool $id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-positions-for-address ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q leveragelp get-positions-for-address $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-positions ..."
$BINARY q leveragelp get-positions

echo "Querying get-status ..."
$BINARY q leveragelp get-status

echo "Querying get-whitelist ..."
$BINARY q leveragelp get-whitelist

echo "Querying user is_whitelisted ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address:"
    read address
    $BINARY q leveragelp is-whitelisted $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying liquidation-price ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the address and position id:"
    read address position_id
    $BINARY q leveragelp liquidation-price $address $position_id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying list-pool ..."
$BINARY q leveragelp list-pool

echo "Querying open-estimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter [amm_pool_id] [collateral] [leverage]"
    read amm_pool_id collateral leverage
    $BINARY q leveragelp open-estimation $amm_pool_id $collateral $leverage
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying rewards ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter address and position id"
    read address position_id
    $BINARY q leveragelp rewards $address $position_id
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying show-pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter index"
    read index
    $BINARY q leveragelp show-pool $index
    echo "Want to make another query, true/false"
    read ASK_USER
done
