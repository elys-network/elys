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

echo "Starting Amm module's Txns"

echo "Tx: Create Pool"
$BINARY tx amm create-pool 100uelys,100uusdc 1000000000uelys,1000000000uusdc --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false  --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000

query_tx "create pool"

echo "Tx: Join the Above created pool, enter pool_id: "
read pool_id
$BINARY tx amm join-pool $pool_id 1200uelys,1200uusdc 10000000 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000
query_tx "Join pool"
$BINARY tx amm join-pool $pool_id 1200uelys,1200uusdc 10000000 --from=$MY_TEST1_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000
query_tx "Join pool"
$BINARY tx amm join-pool $pool_id 1200uelys,1200uusdc 10000000 --from=$MY_TEST2_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000
query_tx "Join pool"
$BINARY tx amm join-pool $pool_id 1200uelys,1200uusdc 10000000 --from=$MY_TEST3_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000
query_tx "Join pool"
$BINARY tx amm join-pool $pool_id 1200uelys,1200uusdc 10000000 --from=$MY_TEST4_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000
query_tx "Join pool"
$BINARY tx amm join-pool $pool_id 1200uelys,1200uusdc 10000000 --from=$MY_TEST5_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000
query_tx "Join pool"


echo "Tx: Exit the Above pool"
$BINARY tx amm exit-pool $pool_id 1uelys,1uusdc 1000000000000 --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --yes --gas=1000000

query_tx "Exit pool"

echo "Tx: Swap exact amount in from the Above pool"
$BINARY tx amm swap-exact-amount-in 100000uusdc 10000 $pool_id uelys --from=my_validator --keyring-backend=test --chain-id=elys --yes --gas=1000000

query_tx "Swap exact amount in"

echo "Tx: Swap exact amount out from the Above pool"
$BINARY tx amm swap-exact-amount-out 10000uelys 20000 $pool_id uusdc --from=my_validator --keyring-backend=test --chain-id=elys --yes --gas=1000000

query_tx "Swap exact amount out"

echo "Tx: Swap by denom from the Above pool"
$BINARY tx amm swap-by-denom 1000000uelys uelys uusdc --min-amount=100000uusdc --max-amount=1000000uelys --from=my_validator --keyring-backend=test --chain-id=elys --yes --gas=1000000

query_tx "Swap by denom"


# CHANGE PARAMS ACCORDINGLY in update_params.json
# echo "Tx: Create and vote on gov Proposal for update pool params"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/amm/update_pool_params.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh

# echo "Tx: Create and vote on gov Proposal for update params"
# $BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/amm/update_params.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
# query_tx "create gov proposal"
# sh ${current_dir}/scripts/test_tx_query/vote.sh