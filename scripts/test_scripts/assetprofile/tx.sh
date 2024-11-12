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

$BINARY tx assetprofile add-entry xyz 18 xyz "" "" "" "XYZ" "" "" "" "" "" [] "" "" "" true true --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "add entry"
$BINARY tx assetprofile add-entry zyza 18 zyza "" "" "" "XYZA" "" "" "" "" "" [] "" "" "" true true --from=$MY_TEST1_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "add entry"
$BINARY tx assetprofile add-entry zyzab 18 zyzab "" "" "" "XYZB" "" "" "" "" "" [] "" "" "" true true --from=$MY_TEST2_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "add entry"
$BINARY tx assetprofile add-entry zyzac 18 zyzac "" "" "" "XYZC" "" "" "" "" "" [] "" "" "" true true --from=$MY_TEST3_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "add entry"
$BINARY tx assetprofile add-entry zyzad 18 zyzad "" "" "" "XYZD" "" "" "" "" "" [] "" "" "" true true --from=$MY_TEST4_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "add entry"
$BINARY tx assetprofile add-entry zyzae 18 zyzae "" "" "" "XYZE" "" "" "" "" "" [] "" "" "" true true --from=$MY_TEST5_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "add entry"

echo "Tx: Create and vote on gov Proposal for update entry"
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/assetprofile/update_entry.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"

sh ${current_dir}/scripts/test_tx_query/vote.sh

sleep 6s

echo "Tx: Create and vote on gov Proposal for Delete entry"
$BINARY tx gov submit-proposal ${current_dir}/scripts/test_tx_query/assetprofile/delete_entry.json --from=$MY_VALIDATOR_ADDRESS --keyring-backend=test --chain-id=elys --gas=1000000
query_tx "create gov proposal"

sh ${current_dir}/scripts/test_tx_query/vote.sh