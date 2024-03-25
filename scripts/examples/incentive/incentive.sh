#!/usr/bin/env bash

# Local test
elysd tx gov submit-proposal proposal.json --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
elysd tx gov vote 1 Yes --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000