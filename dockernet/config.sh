#!/bin/bash

set -eu
DOCKERNET_HOME=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

STATE=$DOCKERNET_HOME/state
LOGS=$DOCKERNET_HOME/logs
UPGRADES=$DOCKERNET_HOME/upgrades
SRC=$DOCKERNET_HOME/src
PEER_PORT=26656
DOCKER_COMPOSE="docker-compose -f $DOCKERNET_HOME/docker-compose.yml"

# Logs
ELYS_LOGS=$LOGS/elys.log
TX_LOGS=$DOCKERNET_HOME/logs/tx.log
KEYS_LOGS=$DOCKERNET_HOME/logs/keys.log

# List of other chains enabled 
OTHER_CHAINS=(BAND)

# Sets up upgrade if {UPGRADE_NAME} is non-empty
UPGRADE_NAME=""
UPGRADE_OLD_COMMIT_HASH=""

# COIN TYPES
# Coin types can be found at https://github.com/satoshilabs/slips/blob/master/slip-0044.md
COSMOS_COIN_TYPE=118
BAND_COIN_TYPE=494
ETH_COIN_TYPE=60
TERRA_COIN_TYPE=330

# CHAIN PARAMS
BLOCK_TIME='1s'
UNBONDING_TIME="120s"
MAX_DEPOSIT_PERIOD="30s"
VOTING_PERIOD="30s"

VAL_TOKENS=5000000000000
STAKE_TOKENS=5000000000
ADMIN_TOKENS=1000000000

# CHAIN MNEMONICS
VAL_MNEMONIC_1="close soup mirror crew erode defy knock trigger gather eyebrow tent farm gym gloom base lemon sleep weekend rich forget diagram hurt prize fly"
VAL_MNEMONIC_2="turkey miss hurry unable embark hospital kangaroo nuclear outside term toy fall buffalo book opinion such moral meadow wing olive camp sad metal banner"
VAL_MNEMONIC_3="tenant neck ask season exist hill churn rice convince shock modify evidence armor track army street stay light program harvest now settle feed wheat"
VAL_MNEMONIC_4="tail forward era width glory magnet knock shiver cup broken turkey upgrade cigar story agent lake transfer misery sustain fragile parrot also air document"
VAL_MNEMONIC_5="crime lumber parrot enforce chimney turtle wing iron scissors jealous indicate peace empty game host protect juice submit motor cause second picture nuclear area"
VAL_MNEMONICS=(
    "$VAL_MNEMONIC_1"
    "$VAL_MNEMONIC_2"
    "$VAL_MNEMONIC_3"
    "$VAL_MNEMONIC_4"
    "$VAL_MNEMONIC_5"
)
REV_MNEMONIC="tonight bonus finish chaos orchard plastic view nurse salad regret pause awake link bacon process core talent whale million hope luggage sauce card weasel"

# ELYS 
ELYS_CHAIN_ID=ELYS
ELYS_NODE_PREFIX=elys
ELYS_NUM_NODES=3
ELYS_VAL_PREFIX=val
ELYS_ADDRESS_PREFIX=elys
ELYS_DENOM=uelys
ELYS_RPC_PORT=26657
ELYS_ADMIN_ACCT=admin

# Binaries are contingent on whether we're doing an upgrade or not
if [[ "$UPGRADE_NAME" == "" ]]; then 
  ELYS_BINARY="$DOCKERNET_HOME/../build/elysd"
else
  if [[ "${NEW_BINARY:-false}" == "false" ]]; then
    ELYS_BINARY="$UPGRADES/binaries/elysd1"
  else
    ELYS_BINARY="$UPGRADES/binaries/elysd2"
  fi
fi
ELYS_MAIN_CMD="$ELYS_BINARY --home $DOCKERNET_HOME/state/${ELYS_NODE_PREFIX}1"

# BAND 
BAND_CHAIN_ID=BAND
BAND_NODE_PREFIX=band
BAND_NUM_NODES=1
BAND_BINARY="$DOCKERNET_HOME/../build/bandd"
BAND_VAL_PREFIX=gval
BAND_REV_ACCT=grev1
BAND_ADDRESS_PREFIX=band
BAND_DENOM=uband
BAND_RPC_PORT=26557
BAND_MAIN_CMD="$BAND_BINARY --home $DOCKERNET_HOME/state/${BAND_NODE_PREFIX}1"

# RELAYER
RELAYER_CMD="$DOCKERNET_HOME/../build/relayer --home $STATE/relayer"
RELAYER_BAND_EXEC="$DOCKER_COMPOSE run --rm relayer-band"

RELAYER_ELYS_ACCT=rly1
RELAYER_BAND_ACCT=rly2
RELAYER_ACCTS=($RELAYER_BAND_ACCT)

RELAYER_BAND_MNEMONIC="fiction perfect rapid steel bundle giant blade grain eagle wing cannon fever must humble dance kitchen lazy episode museum faith off notable rate flavor"
RELAYER_MNEMONICS=(
  "$RELAYER_BAND_MNEMONIC"
)

PRICE_FEEDER_ELYS_ACCT=feeder
PRICE_FEEDER_ACCTS=($PRICE_FEEDER_ELYS_ACCT)
PRICE_FEEDER_MNEMONIC="surprise pear reject sail fresh west equal tragic rain divorce direct roast piece canyon rival amateur leaf poet admit frequent point measure vapor abstract"
PRICE_FEEDER_MNEMONICS=(
  "$PRICE_FEEDER_MNEMONIC"
)

ELYS_ADDRESS() { 
  # After an upgrade, the keys query can sometimes print migration info, 
  # so we need to filter by valid addresses using the prefix
  $ELYS_MAIN_CMD keys show ${ELYS_VAL_PREFIX}1 --keyring-backend test -a | grep $ELYS_ADDRESS_PREFIX
}
BAND_ADDRESS() { 
  $BAND_MAIN_CMD keys show ${BAND_VAL_PREFIX}1 --keyring-backend test -a | grep $BAND_ADDRESS_PREFIX
}

CSLEEP() {
  for i in $(seq $1); do
    sleep 1
    printf "\r\t$(($1 - $i))s left..."
  done
}

GET_VAR_VALUE() {
  var_name="$1"
  echo "${!var_name}"
}

WAIT_FOR_BLOCK() {
  num_blocks="${2:-1}"
  for i in $(seq $num_blocks); do
    ( tail -f -n0 $1 & ) | grep -q "INF executed block height="
  done
}

WAIT_FOR_STRING() {
  ( tail -f -n0 $1 & ) | grep -q "$2"
}

WAIT_FOR_BALANCE_CHANGE() {
  chain=$1
  address=$2
  denom=$3

  max_blocks=30

  main_cmd=$(GET_VAR_VALUE ${chain}_MAIN_CMD)
  initial_balance=$($main_cmd q bank balances $address --denom $denom | grep amount)
  for i in $(seq $max_blocks); do
    new_balance=$($main_cmd q bank balances $address --denom $denom | grep amount)

    if [[ "$new_balance" != "$initial_balance" ]]; then
      break
    fi

    WAIT_FOR_BLOCK $ELYS_LOGS 1
  done
}

GET_VAL_ADDR() {
  chain=$1
  val_index=$2

  MAIN_CMD=$(GET_VAR_VALUE ${chain}_MAIN_CMD)
  $MAIN_CMD q staking validators | grep ${chain}_${val_index} -A 5 | grep operator | awk '{print $2}'
}

TRIM_TX() {
  grep -E "code:|txhash:" | sed 's/^/  /'
}