#!/bin/bash

set -e

BINARY="elysd"
NODE="tcp://localhost:26657"
MY_VALIDATOR_ADDRESS=$($BINARY keys show my_validator -a --keyring-backend test)
BROKER_ADDRESS=$($BINARY keys show test1 -a --keyring-backend test)
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

echo "setting up pools and bonding tokens: "
echo "Tx: Create Pool"
$BINARY tx amm create-pool 100uatom,100uusdc 2240000000000uusdc,180000000000uatom --use-oracle=true --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create pool"

echo "Tx: stablestake Bond, Put funds on stablestake"
$BINARY tx stablestake bond 100000000 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "stablestake Bond"

# # Make sure this is the first pool(pool id is 1) otherwise update proposal.json
# echo "Tx: Create and vote on gov Proposal for Adding pool in leveragelp"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/leveragelp/proposal_add_pool.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# echo "Wait for voting period to end passed"

# sleep 80s

$BINARY tx perpetual open long 5 2 uatom 100000000uusdc --take-profit 45 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open long"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_VALIDATOR_ADDRESS

# sleep 6s

$BINARY tx perpetual open long 5 2 uatom 100000000uusdc --take-profit 45 --from=$MY_TEST2_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open long"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST2_ADDRESS

$BINARY tx perpetual open long 5 2 uatom 100000000uusdc --take-profit 45 --from=$MY_TEST3_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open long"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST3_ADDRESS

$BINARY tx perpetual open long 5 2 uatom 100000000uusdc --take-profit 45 --from=$MY_TEST4_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open long"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST4_ADDRESS

$BINARY tx perpetual open long 5 2 uatom 100000000uusdc --take-profit 45 --from=$MY_TEST5_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open long"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST5_ADDRESS

# echo "Open long position using broker"
# $BINARY tx perpetual broker-open long 5 5 uatom 100000000uusdc $MY_TEST2_ADDRESS --take-profit 45 --from=$BROKER_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "perpetual broker-open long"
# echo "Querying position"
# $BINARY q perpetual get-positions-for-address $MY_TEST2_ADDRESS

# echo "Close position using broker"
# $BINARY tx perpetual broker-close 4 49934694 $MY_TEST2_ADDRESS --from=$BROKER_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "perpetual broker-close"
# echo "Querying position"
# $BINARY q perpetual get-positions-for-address $MY_TEST2_ADDRESS


# sleep 6s

$BINARY tx perpetual open short 5 3 uatom 10000000uusdc --take-profit 1 --stop-loss 10.2 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open short"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_VALIDATOR_ADDRESS

$BINARY tx perpetual open short 5 3 uatom 100000000uusdc --take-profit 1 --stop-loss 10.2 --from=$MY_TEST2_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open short"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST2_ADDRESS

$BINARY tx perpetual open short 5 3 uatom 100000000uusdc --take-profit 1 --stop-loss 10.2 --from=$MY_TEST3_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open short"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST3_ADDRESS

$BINARY tx perpetual open short 5 3 uatom 100000000uusdc --take-profit 1 --stop-loss 10.2 --from=$MY_TEST4_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open short"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST4_ADDRESS

$BINARY tx perpetual open short 5 3 uatom 100000000uusdc --take-profit 1 --stop-loss 10.2 --from=$MY_TEST5_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual open short"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_TEST5_ADDRESS

sleep 6s
$BINARY tx perpetual update-take-profit-price 40 1 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual update-take-profit-price"

$BINARY tx perpetual update-stop-loss 9 1 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "perpetual update-take-profit-price"
echo "Querying position"
$BINARY q perpetual get-positions-for-address $MY_VALIDATOR_ADDRESS

# echo "Tx: Create and vote on gov Proposal for Whitelisting address"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/perpetual/whitelist_proposal.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 90s
# $BINARY q perpetual get-whitelist
# $BINARY q perpetual is-whitelisted elys1679c7fc2sxedkjwfu24qmdax9e86fx6tdqmlj5
# sleep 5s

# echo "Tx: Create and vote on gov Proposal for DEWhitelisting address"
# $BINARY tx gov submit-proposal ${currentf_dir}/scripts/test_tx_query/perpetual/dewhitelist_proposal.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# sleep 90s
# $BINARY q perpetual get-whitelist
# $BINARY q perpetual is-whitelisted elys1679c7fc2sxedkjwfu24qmdax9e86fx6tdqmlj5

# echo "Tx: Create and vote on gov Proposal for params update"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/perpetual/update_params_proposal.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh


# echo "Tx: close-positions"
# $BINARY tx leveragelp close-positions ${current_dir}/scripts/test_tx_query/leveragelp/liquidate.json ${current_dir}/scripts/test_tx_query/leveragelp/stoploss.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
