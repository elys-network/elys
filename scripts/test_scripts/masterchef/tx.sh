#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
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

# echo "Tx: Create and vote on gov Proposal for add_external_reward_denom "
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/masterchef/add_external_reward_denom.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 80s
# $BINARY tx masterchef add-external-incentive uatom 2 3500 3600 10 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# $BINARY q masterchef external-incentive uatom


# $BINARY tx masterchef claim-rewards --pool-ids="1,2" --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000


# echo "Tx: Create and vote on gov Proposal for toggle_pool_eden_rewards "
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/masterchef/toggle_pool_eden_rewards.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

echo "Tx: Create and vote on gov Proposal for update pool multipliers "
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/masterchef/update_pool_multipliers.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"
sh ${current_dir}/scripts/test_tx_query/vote.sh


echo "Tx: Create and vote on gov Proposal for update params "
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/masterchef/update_params.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"
sh ${current_dir}/scripts/test_tx_query/vote.sh