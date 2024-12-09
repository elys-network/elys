#!/usr/bin/env bash

TREASURY=$(elysd keys show -a treasury --keyring-backend=test)

elysd tx amm create-pool elysd tx amm create-pool 10uatom,10uusdt 10000uatom,10000uusdt --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false --from=treasury --yes --gas=1000000

# single asset add-liquidity
elysd tx amm join-pool 1 2000uatom 90000000000000000 true --from=treasury --yes --gas=1000000
# multiple asset add-liquidity
elysd tx amm join-pool 1 2000uatom,2000uusdt 200000000000000000 true --from=treasury --yes --gas=1000000

# swap
elysd tx amm swap-exact-amount-in 10uatom 1 1 uusdt --from=treasury --yes --gas=1000000

elysd query commitment show-commitments $TREASURY
elysd query bank balances $TREASURY
elysd query amm show-pool 1

# query amm swap estimation 5000 ATOM
ATOM=ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4
USDC=ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65
elysd q amm swap-estimation-by-denom 5000000000$ATOM $ATOM $USDC
elysd q amm swap-estimation-by-denom 40000000000$USDC $USDC $ATOM
elysd q amm swap-estimation-by-denom 4000000000uelys uelys $USDC
elysd q amm swap-estimation-by-denom 10000000000$USDC $USDC uelys
elysd q amm swap-estimation-by-denom 5000000000$ATOM $ATOM $USDC --node=https://rpc.testnet.elys.network:443

elysd query amm show-pool 2
