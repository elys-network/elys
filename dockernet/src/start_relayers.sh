#!/bin/bash

set -eu 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/../config.sh

for chain in ${OTHER_CHAINS[@]}; do
    relayer_exec=$(GET_VAR_VALUE RELAYER_${chain}_EXEC)
    chain_name=$(printf "$chain" | awk '{ print tolower($0) }')
    account_name=$(GET_VAR_VALUE RELAYER_${chain}_ACCT)
    mnemonic=$(GET_VAR_VALUE     RELAYER_${chain}_MNEMONIC)
    coin_type=$(GET_VAR_VALUE    ${chain}_COIN_TYPE)

    relayer_logs=${LOGS}/relayer-${chain_name}.log
    relayer_config=$STATE/relayer-${chain_name}/config

    mkdir -p $relayer_config
    chmod -R 777 $STATE/relayer-${chain_name}
    cp ${DOCKERNET_HOME}/config/relayer_config.yaml $relayer_config/config.yaml

    printf "ELYS <> $chain - Adding relayer keys..."
    $relayer_exec rly keys restore elys $RELAYER_ELYS_ACCT "$mnemonic" >> $relayer_logs 2>&1
    $relayer_exec rly keys restore $chain_name $account_name "$mnemonic" --coin-type $coin_type >> $relayer_logs 2>&1
    echo "Done restoring relayer keys"

    printf "ELYS <> $chain - Creating client, connection, and transfer channel..." | tee -a $relayer_logs
    $relayer_exec rly transact link elys-${chain_name} >> $relayer_logs 2>&1
    $relayer_exec rly transact channel elys-${chain_name} --src-port oracle --dst-port oracle --version "bandchain-1" >> $relayer_logs 2>&1
    echo "Done."

    $DOCKER_COMPOSE up -d relayer-${chain_name}
    $DOCKER_COMPOSE logs -f relayer-${chain_name} | sed -r -u "s/\x1B\[([0-9]{1,3}(;[0-9]{1,2})?)?[mGK]//g" >> $relayer_logs 2>&1 &
done
