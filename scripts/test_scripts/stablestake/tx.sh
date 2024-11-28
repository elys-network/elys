#!/bin/bash

# set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
MY_TEST1_ADDRESS=$($BINARY keys show test1 -a --keyring-backend test)
MY_TEST2_ADDRESS=$($BINARY keys show test2 -a --keyring-backend test)
MY_TEST3_ADDRESS=$($BINARY keys show test3 -a --keyring-backend test)
MY_TEST4_ADDRESS=$($BINARY keys show test4 -a --keyring-backend test)
MY_TEST5_ADDRESS=$($BINARY keys show test5 -a --keyring-backend test)
current_dir=$(pwd)

query_tx() {
    local tx_type=$1
    echo "Query the ${tx_type} txhash, enter txhash:"
    read -r tx_hash
    $BINARY q tx "$tx_hash" --node "$NODE"
}

echo "Tx: stablestake Bond"
$BINARY tx stablestake bond 100000000 --from=$MY_TEST1_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"
echo "Tx: stablestake Bond"
$BINARY tx stablestake bond 100000000 --from=$MY_TEST2_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"
echo "Tx: stablestake Bond"
$BINARY tx stablestake bond 100000000 --from=$MY_TEST3_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"
echo "Tx: stablestake Bond"
$BINARY tx stablestake bond 100000000 --from=$MY_TEST4_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"
echo "Tx: stablestake Bond"
$BINARY tx stablestake bond 100000000 --from=$MY_TEST5_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"

sleep 4s

echo "Tx: stablestake UnBond"
$BINARY tx stablestake unbond 10000000 --from=$MY_TEST1_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake UnBond"
echo "Tx: stablestake UnBond"
$BINARY tx stablestake unbond 10000000 --from=$MY_TEST2_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake UnBond"
echo "Tx: stablestake UnBond"
$BINARY tx stablestake unbond 10000000 --from=$MY_TEST3_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake UnBond"
echo "Tx: stablestake UnBond"
$BINARY tx stablestake unbond 10000000 --from=$MY_TEST4_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake UnBond"
echo "Tx: stablestake UnBond"
$BINARY tx stablestake unbond 10000000 --from=$MY_TEST5_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake UnBond"


# echo "Tx: Create and vote on gov Proposal for Params update"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/stablestake/update_params.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh
