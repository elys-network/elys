#!/usr/bin/env bash

TREASURY=$(elysd keys show -a treasury --keyring-backend=test)

elysd tx amm create-pool 10uatom,10uusdt 10000uatom,10000uusdt 0.00 0.00 --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# single asset add-liquidity
elysd tx amm join-pool 0 2000uatom 90000000000000000 true --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
# multiple asset add-liquidity
elysd tx amm join-pool 0 2000uatom,2000uusdt 200000000000000000 true --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# swap
elysd tx amm swap-exact-amount-in 10uatom 1 0 uusdt --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

elysd query commitment show-commitments $TREASURY
elysd query bank balances $TREASURY
elysd query amm show-pool 0