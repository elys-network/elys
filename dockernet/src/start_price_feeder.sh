#!/bin/bash

set -eu 
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source ${SCRIPT_DIR}/../config.sh

price_feeder_logs=${LOGS}/price-feeder.log
price_feeder_config=$STATE/price-feeder
mkdir -p $price_feeder_config
chmod -R 777 $price_feeder_config
cp ${DOCKERNET_HOME}/config/price-feeder.toml $price_feeder_config/price-feeder.toml

$DOCKER_COMPOSE up -d price-feeder
$DOCKER_COMPOSE logs -f price-feeder | sed -r -u "s/\x1B\[([0-9]{1,3}(;[0-9]{1,2})?)?[mGK]//g" >> $price_feeder_logs 2>&1 &
