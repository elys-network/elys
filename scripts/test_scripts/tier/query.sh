#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running tier module's query"

echo "Querying Params ..."
$BINARY q tier params --node $NODE

echo "Querying list-portfolio ..."
$BINARY q tier list-portfolio --node $NODE

echo "Querying calculate-discount ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] of the user:"
    read address
    $BINARY q tier calculate-discount $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-amm-price ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [denom] [decimal]:"
    read denom decimal
    $BINARY q tier get-amm-price $denom $decimal
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying get-consolidated-price ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [denom] :"
    read denom
    $BINARY q tier get-consolidated-price $denom
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying leverage-lp-total ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier leverage-lp-total $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying liquid-total ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier liquid-total $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying perpetual ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier perpetual $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying rewards-total ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier rewards-total $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying show-portfolio ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier show-portfolio $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying staked ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier staked $address
    echo "Want to make another query, true/false"
    read ASK_USER
done

echo "Querying staked-pool ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the [address] :"
    read address
    $BINARY q tier staked-pool $address
    echo "Want to make another query, true/false"
    read ASK_USER
done


