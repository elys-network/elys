#!/usr/bin/env bash

elysd query commitment number-of-commitments
# number: "395166"

elysd tx commitment commit-claimed-rewards 503544 ueden --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
elysd tx commitment commit-claimed-rewards 1678547 uedenb --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
