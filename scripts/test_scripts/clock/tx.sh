#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
MY_TEST1_ADDRESS=$($BINARY keys show test1 -a --keyring-backend test)
MY_TEST2_ADDRESS=$($BINARY keys show test2 -a --keyring-backend test)
current_dir=$(pwd)
VALIDATOR_OPERATOR_ADDRESS={$VALIDATOR_OPERATOR_ADDRESS:="elysvaloper1j96kg6nq4l00q3qcpjzf093f2xsn3lvksezmpw"}

query_tx() {
    local tx_type=$1
    echo "Query the ${tx_type} txhash, enter txhash:"
    read -r tx_hash
    $BINARY q tx "$tx_hash" --node "$NODE"
}

echo "Tx: Create and vote on gov Proposal for updating params"
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/clock/update_params.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"
sh ${current_dir}/scripts/test_tx_query/vote.sh

sleep 90s

sh ${current_dir}/scripts/test_tx_query/clock/query.sh