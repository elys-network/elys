#!/bin/bash

set -eu 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
source ${SCRIPT_DIR}/config.sh

# cleanup any stale state
make stop-docker
rm -rf $STATE $LOGS 
mkdir -p $STATE
mkdir -p $LOGS


# If we're testing an upgrade, setup cosmovisor
if [[ "$UPGRADE_NAME" != "" ]]; then
    printf "\n>>> UPGRADE ENABLED! ($UPGRADE_NAME)\n\n"
    
    # Update binary #2 with the binary that was just compiled
    mkdir -p $UPGRADES/binaries
    rm -f $UPGRADES/binaries/elysd2
    cp $DOCKERNET_HOME/../build/elysd $UPGRADES/binaries/elysd2

    # Build a cosmovisor image with the old binary and replace the elys docker image with a new one
    #  that has both binaries and is running cosmovisor
    # The reason for having a separate cosmovisor image is so we can cache the building of cosmovisor and the old binary
    echo "Building Cosmovisor..."
    docker build \
        -t elyszone:cosmovisor \
        --build-arg old_commit_hash=$UPGRADE_OLD_COMMIT_HASH \
        -f $UPGRADES/Dockerfile.cosmovisor .

    echo "Done"
fi


# Initialize the state for each chain
for chain in ELYS ${OTHER_CHAINS[@]}; do
    bash $SRC/init_chain.sh $chain
done


# Start the chain and create the transfer channels
bash $SRC/start_chain.sh 
bash $SRC/start_relayers.sh 

# Register all host zones 
for i in ${!OTHER_CHAINS[@]}; do
    bash $SRC/register_host.sh ${OTHER_CHAINS[$i]} $i 
done

$SRC/create_logs.sh &
