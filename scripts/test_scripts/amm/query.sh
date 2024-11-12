#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Amm module's query"

echo "Querying Params ..."
$BINARY q oracle params --node $NODE

echo "Querying PoolAll..."
$BINARY q amm list-pool --node $NODE

echo "Querying Pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the pool-id:"
    read POOL_ID
    $BINARY q amm show-pool $POOL_ID --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying List DenomLiquidity..."
$BINARY q amm list-denom-liquidity --node $NODE

echo "Querying DenomLiquidity ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the denom:"
    read DENOM
    $BINARY q amm show-denom-liquidity $DENOM --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

Handle array input
echo "Querying SwapEstimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter token_in:"
    read token_in
    echo "Enter SwapAmountIn routes:"
    read -a routes
    $BINARY q amm swap-estimation $token_in $routes --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying JoinPoolEstimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter pool-id and token_in (space separated):"
    read pool_id token_in
    $BINARY q amm join-pool-estimation $pool_id $token_in --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying ExitPoolEstimation ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter pool-id, token_in token_out_denom (space separated):"
    read pool_id token_in token_out_denom
    $BINARY q amm exit-pool-estimation $pool_id $token_in $token_out_denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying SlippageTrackAll..."
$BINARY q amm tracked-slippage-all --node $NODE

echo "Querying SlippageTrack ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter pool-id:"
    read pool_id
    $BINARY q amm tracked-slippage $pool_id --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying Balance ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter address and denom (space separated):"
    read address denom
    $BINARY q amm balance $address $denom --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying InRouteByDenom ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter denom-in and denom-out (space separated):"
    read denom_in denom_out
    $BINARY q amm in-route-by-denom $denom_in $denom_out --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying OutRouteByDenom ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter denom-out and denom-in (space separated):"
    read denom_out denom_in
    $BINARY q amm out-route-by-denom $denom_out $denom_in --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying SwapEstimationByDenom ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter amount, denom-in and denom-out (space separated):"
    read amount denom_in denom_out
    $BINARY q amm swap-estimation-by-denom $amount $denom_in $denom_out --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done