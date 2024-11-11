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

echo "Tx: Stake"
$BINARY tx commitment stake 200 uedenb elysvaloper16sm9uuk59wg5fzvd7d8u9erwaj44x7f9f90eqz  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "Stake uedenb"

$BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS


echo "Tx: UnStake"
$BINARY tx commitment unstake 200 uedenb elysvaloper16sm9uuk59wg5fzvd7d8u9erwaj44x7f9f90eqz  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "UnStake uedenb"

$BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS

echo "Tx: uncommit-tokens"
$BINARY tx commitment uncommit-tokens 200 uedenb  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "uncommit-tokens"

$BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS

echo "Tx: commit-claimed-rewards"
$BINARY tx commitment commit-claimed-rewards 100 uedenb  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "uncommit-tokens"

$BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS

echo "Tx: Create and vote on gov Proposal for update_vesting_info"
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/commitment/update_vesting_info.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"
sh ${current_dir}/scripts/test_tx_query/vote.sh

sleep 90s

echo "Tx: vest"
$BINARY tx commitment vest 98 uedenb  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "vest"

$BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS

echo "Tx: vest-liquid"
$BINARY tx commitment vest-liquid 4560 uelys  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "vest-liquid"

$BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS

# NOT ENABLED YET FROM THE CODE
# echo "Tx: vest-now"
# $BINARY tx commitment vest-now 10 uedenb  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "vest-now"

# $BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS

# echo "Tx: cancel-vest"
# $BINARY tx commitment cancel-vest 98 ueden --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "cancel-vest"

# $BINARY q commitment show-commitments $MY_VALIDATOR_ADDRESS
