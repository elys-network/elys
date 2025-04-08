#!/usr/bin/env bash

# query params
elysd query bank balances $(elysd keys show -a treasury --keyring-backend=test)

elysd query stablestake params
# params:
#   deposit_denom: uusdc
#   epoch_length: "1"
#   health_gain_factor: "1.000000000000000000"
#   interest_rate: "0.150000000000000000"
#   interest_rate_decrease: "0.010000000000000000"
#   interest_rate_increase: "0.010000000000000000"
#   interest_rate_max: "0.170000000000000000"
#   interest_rate_min: "0.120000000000000000"
#   redemption_rate: "1.000000000000000000"
#   total_value: "100000000"
elysd query stablestake borrow-ratio
# borrow_ratio: "0.000000000000000000"
# total_borrow: "0"
# total_deposit: "100000000"

elysd tx stablestake bond 100000000 --from=treasury --yes --gas=1000000


# Testnet
elysd tx stablestake bond 112131 --from=t2a --yes --gas=1000000 --node=https://rpc.testnet.elys.network:443 --fees=250uelys
