#!/usr/bin/env bash

# queries
elysd query launchpad all-orders
elysd query launchpad bonus $(elysd keys show -a treasury --keyring-backend=test)
elysd query launchpad buy-elys-estimation uusdc 1000000
elysd query launchpad return-elys-estimation 1 1000000
elysd query launchpad orders $(elysd keys show -a treasury --keyring-backend=test)
elysd query launchpad params
elysd query launchpad module-balances

# Buy Elys
elysd tx launchpad buy-elys uatom 100000000000 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# Return Elys
elysd tx launchpad return-elys 1 100000000000 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# Deposit Elys token
elysd tx launchpad deposit-elys-token 100000000000uelys --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# Withdraw raised
elysd tx launchpad withdraw-raised 100000000000uatom --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
