#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"

echo "Running Oracle module's query"

echo "Querying Params ..."
$BINARY q oracle params --node $NODE

echo "Querying LastBandRequestId ..."
$BINARY q oracle last-band-request-id --node $NODE

echo "Querying BandPriceResult ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the request-id:"
    read REQUEST_ID
    $BINARY q oracle band-price-result $REQUEST_ID --node $NODE
    echo "Want to make another query, true/false"
    read ASK_USER
done


echo "Querying AssetInfoAll ..."
$BINARY q oracle list-asset-info --node $NODE

echo "Querying AssetInfo ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter the asset's denom to query the AssetInfo:"
    read DENOM
    $BINARY q oracle show-asset-info $DENOM --node $NODE
    echo "Want to query other asset-info, true/false"
    read ASK_USER
done

echo "Querying all price ..."
$BINARY q oracle list-price --node $NODE

echo "Querying Price ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter asset to query its Price:"
    read ASSET
    $BINARY q oracle show-price $ASSET --node $NODE
    echo "Want to query other asset's price, true/false"
    read ASK_USER
done

echo "Querying all price feeder ..."
$BINARY q oracle list-price-feeder --node $NODE

echo "Querying Price Feeder ..."
ASK_USER=true
while [ "$ASK_USER" = true ]; do
    echo "Enter price feeder address:"
    read FEEDER
    $BINARY q oracle show-price-feeder $FEEDER --node $NODE
    echo "Want to query other price feeder, true/false"
    read ASK_USER
done
