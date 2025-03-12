#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show validator -a --keyring-backend test)
# MY_TEST1_ADDRESS=$($BINARY keys show test1 -a --keyring-backend test)
# MY_TEST2_ADDRESS=$($BINARY keys show test2 -a --keyring-backend test)
current_dir=$(pwd)

query_tx() {
    local tx_type=$1
    echo "Query the ${tx_type} txhash, enter txhash:"
    read -r tx_hash
    $BINARY q tx "$tx_hash" --node "$NODE"
}

echo "Voting: Enter gov proposal id: "
read proposal_id
# response=$(curl -s http://localhost:1317/cosmos/gov/v1beta1/proposals)
# proposal_id=$(echo "$response" | jq -r '.pagination.total')
$BINARY tx gov vote $proposal_id Yes --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elysicstestnet-1 --gas=1000000
query_tx "gov vote"