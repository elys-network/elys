#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
BROKER_ADDRESS=$($BINARY keys show test1 -a --keyring-backend test)
MY_TEST2_ADDRESS=$($BINARY keys show test2 -a --keyring-backend test)
current_dir=$(pwd)

query_tx() {
    local tx_type=$1
    echo "Query the ${tx_type} txhash, enter txhash:"
    read -r tx_hash
    $BINARY q tx "$tx_hash" --node "$NODE"
}

echo "Tx: Create and vote on gov Proposal for MsgUpdateMinCommission"
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/parameter/update_min_commission.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"
sh ${current_dir}/scripts/test_tx_query/vote.sh
sleep 6s
sh ${current_dir}/scripts/test_tx_query/parameter/query.sh

# echo "Tx: Create and vote on gov Proposal for MsgUpdateMaxVotingPower"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/parameter/update_max_voting_power.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 10s

# echo "Tx: Create and vote on gov Proposal for MsgUpdateMinSelfDelegation"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/parameter/update_min_self_delegation.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 10s

# echo "Tx: Create and vote on gov Proposal for MsgUpdateBrokerAddress"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/parameter/update_broker_address.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 10s

# echo "Tx: Create and vote on gov Proposal for MsgUpdateTotalBlocksPerYear"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/parameter/update_total_block_per_year.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 10s

# echo "Tx: Create and vote on gov Proposal for MsgUpdateRewardsDataLifetime"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/parameter/update_rewards_data_lifetime.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh
