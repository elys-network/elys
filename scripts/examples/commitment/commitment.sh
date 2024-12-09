#!/usr/bin/env bash

elysd query commitment number-of-commitments
# number: "395166"

elysd query commitment committed-tokens-locked $(elysd keys show -a treasury --keyring-backend=test)
# address: elys12tzylat4udvjj56uuhu3vj2n4vgp7cf9fwna9w
# locked_committed:
# - amount: "110000000000000000000000"
#   denom: amm/pool/1
# total_committed:
# - amount: "110000000000000000000000"
#   denom: amm/pool/1
# - amount: "100000000"
#   denom: stablestake/share

elysd tx commitment commit-claimed-rewards 503544 ueden --from=treasury --yes --gas=1000000
elysd tx commitment commit-claimed-rewards 1678547 uedenb --from=treasury --yes --gas=1000000
