#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
MY_TEST2_ADDRESS=$($BINARY keys show test2 -a --keyring-backend test)
current_dir=$(pwd)

query_tx() {
    local tx_type=$1
    echo "Query the ${tx_type} txhash, enter txhash:"
    read -r tx_hash
    $BINARY q tx "$tx_hash" --node "$NODE"
}

echo "Tx: Create Asset info"
$BINARY tx oracle create-asset-info udum DUM DUM DUM 18 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create create-asset-info"

sleep 3s

$BINARY q oracle show-asset-info udum

# sleep 3s

echo "Tx: feed-price"
$BINARY tx oracle feed-price udum 12.0 elys --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "feed-price"

echo "Tx: feed-price"
$BINARY tx oracle feed-price udumtwo 13.0 elys --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "feed-price"

echo "Tx: feed-price"
$BINARY tx oracle feed-price udumqqqtwo 13.0 elys --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "feed-price"

echo "Tx: feed-price"
$BINARY tx oracle feed-price udumwwwwwtwo 13.0 elys --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "feed-price"

sleep 3s

echo "Tx: set-price-feeder InActive"
$BINARY tx oracle set-price-feeder false --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-price-feeder"

$BINARY q oracle show-price-feeder $MY_VALIDATOR_ADDRESS

echo "Tx: set-price-feeder Active"
$BINARY tx oracle set-price-feeder true --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-price-feeder"

$BINARY q oracle show-price-feeder $MY_VALIDATOR_ADDRESS

# GOVERNANCE MODULE's TXN
# 1. Update params
# 2. Add price feeders
# 3. Remove price feeders
# 4. Remove assetinfo

# echo "Tx: Create and vote on gov Proposal for Params update"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/oracle/update_params.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# echo "Tx: Create and vote on gov Proposal for Adding Price Feeders in oracle module"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/oracle/add_price_feeder.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# echo "Tx: Create and vote on gov Proposal for Removing Price Feeders in oracle module"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/oracle/remove_price_feeders.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# echo "Tx: Create and vote on gov Proposal for Removing Asset Info in oracle module"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/oracle/remove_asset_info.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 80s

$BINARY q oracle params
$BINARY q oracle list-price-feeder
$BINARY q oracle list-asset-info
