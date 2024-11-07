#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
BROKER_ADDRESS=$($BINARY keys show test1 -a --keyring-backend test)
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

echo "Tx: set-portfolio"
$BINARY tx tier set-portfolio $MY_TEST1_ADDRESS --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-portfolios"

echo "Tx: set-portfolio"
$BINARY tx tier set-portfolio $MY_TEST2_ADDRESS --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-portfolios"

echo "Tx: set-portfolio"
$BINARY tx tier set-portfolio $MY_TEST3_ADDRESS --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-portfolios"

echo "Tx: set-portfolio"
$BINARY tx tier set-portfolio $MY_TEST4_ADDRESS --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-portfolios"

echo "Tx: set-portfolio"
$BINARY tx tier set-portfolio $MY_TEST5_ADDRESS --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "set-portfolios"

echo "Show portfolio of user test1"
$BINARY q tier show-portfolio $MY_TEST1_ADDRESS 
echo "Show portfolio of user test2"
$BINARY q tier show-portfolio $MY_TEST2_ADDRESS 
echo "Show portfolio of user test3"
$BINARY q tier show-portfolio $MY_TEST3_ADDRESS 
echo "Show portfolio of user test4"
$BINARY q tier show-portfolio $MY_TEST4_ADDRESS 
echo "Show portfolio of user test5"
$BINARY q tier show-portfolio $MY_TEST5_ADDRESS 


