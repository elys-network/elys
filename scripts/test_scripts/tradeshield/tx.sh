#!/bin/bash

set -e

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

# Create spot order
$BINARY tx tradeshield create-spot-order stoploss 1000000uusdc uatom 50 --from=my_validator --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create spot order"
$BINARY tx tradeshield create-perpetual-open-order long 2 1 uatom 1000000uusdc 10 --take-profit 45 --stop-loss 2 --from=my_validator --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create perpetual order"

# $BINARY tx tradeshield execute-orders "1" "1" --from=my_validator --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "execute orders"

# $BINARY q tradeshield list-pending-spot-order

