#!/bin/bash

set -e

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

echo "Starting Leveragelp module's Txns"

# echo "Feed price in oracle"
# $BINARY tx oracle feed-price uatom 4 local --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000

echo "Tx: Create Pool"
$BINARY tx amm create-pool 10uatom,10uusdc 10000000000uatom,10000000000uusdc --use-oracle=true --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create pool"

echo "Tx: stablestake Bond, Put funds on stablestake"
$BINARY tx stablestake bond 100000000 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"

sleep 6s

# Make sure this is the first pool(pool id is 1) otherwise update proposal.json
echo "Tx: Create and vote on gov Proposal for Adding pool in leveragelp"
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/leveragelp/proposal_add_pool.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"

sh ${current_dir}/scripts/test_tx_query/vote.sh

$BINARY query gov proposals

sleep 80s

echo "Tx: Open position leveragelp"
$BINARY tx leveragelp open 5.0 uusdc 5000000 2 0.0 --from=$MY_VALIDATOR_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "open leveragelp position"
$BINARY tx leveragelp open 5.0 uusdc 5000000 2 0.0 --from=$MY_TEST1_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "open leveragelp position"
$BINARY tx leveragelp open 5.0 uusdc 5000000 2 0.0 --from=$MY_TEST2_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "open leveragelp position"
$BINARY tx leveragelp open 5.0 uusdc 5000000 2 0.0 --from=$MY_TEST3_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "open leveragelp position"
$BINARY tx leveragelp open 5.0 uusdc 5000000 2 0.0 --from=$MY_TEST4_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "open leveragelp position"
$BINARY tx leveragelp open 5.0 uusdc 5000000 2 0.0 --from=$MY_TEST5_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "open leveragelp position"

# Wait for 1 hour before withdrawing the token
# echo "Tx: Close position leveragelp"
# $BINARY tx leveragelp close 1 49500000000000000000 --from=$MY_VALIDATOR_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
# query_tx "Close leveragelp position"

echo "Query leveragelp rewards for $MY_VALIDATOR_ADDRESS"
$BINARY query leveragelp rewards $MY_VALIDATOR_ADDRESS 1 --output=json

echo "Query leveragelp rewards for $MY_TEST2_ADDRESS"
$BINARY query leveragelp rewards $MY_TEST2_ADDRESS 1 --output=json

echo "Tx: Claim leveragelp rewards"
$BINARY tx leveragelp claim-rewards 1 --from=$MY_VALIDATOR_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "Claim leveragelp rewards"

echo "Tx: Update Stop loss price"
$BINARY tx leveragelp update-stop-loss 2 10.0 --from=$MY_VALIDATOR_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "Update Stop loss price"

echo "Tx: Update Stop loss price"
$BINARY tx leveragelp update-stop-loss 3 10.0 --from=$MY_TEST2_ADDRESS --chain-id elys --gas 5000000 --keyring-backend="test"
query_tx "Update Stop loss price"


# # Update params with proposal
# echo "Tx: Create and vote on gov Proposal for update params in leveragelp"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/leveragelp/update_params_proposal.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"

# echo "Voting: Enter gov proposal id: "
# read proposal_id
# $BINARY tx gov vote $proposal_id Yes --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "gov vote"

# # Remove pool with proposal
# echo "Tx: Create and vote on gov Proposal o remove pool in leveragelp"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/leveragelp/proposal_remove_pool.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"

# echo "Voting: Enter gov proposal id: "
# read proposal_id
# $BINARY tx gov vote $proposal_id Yes --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "gov vote"

# echo "Tx: Create and vote on gov Proposal for whitelist"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/leveragelp/whitelist_proposal.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 90s

# echo "Tx: Create and vote on gov Proposal for Dewhitelist"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/leveragelp/dewhitelist_proposal.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# echo "Tx: close-positions"
# $BINARY tx leveragelp close-positions ${current_dir}/scripts/test_tx_query/leveragelp/liquidate.json ${current_dir}/scripts/test_tx_query/leveragelp/stoploss.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000