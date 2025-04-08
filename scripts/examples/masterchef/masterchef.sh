#!/usr/bin/env bash

TREASURY=$(elysd keys show -a treasury --keyring-backend=test)

elysd tx amm create-pool 10uatom,10uusdc 20000000000uatom,20000000000uusdc --swap-fee=0.01 --exit-fee=0.00 --use-oracle=false --from=treasury --yes --gas=1000000

# single asset add-liquidity
elysd tx amm join-pool 1 2000uatom 90000000000000000 --from=treasury --yes --gas=1000000
# multiple asset add-liquidity
elysd tx amm join-pool 1 20000000000uatom,20000000000uusdt 1000000000000000000 --from=treasury --yes --gas=1000000

# swap
elysd tx amm swap-exact-amount-in 10000000uatom 1 1 uusdc --from=treasury --yes --gas=1000000

elysd query commitment show-commitments $TREASURY
elysd query bank balances $TREASURY

elysd query masterchef pool-info 1
elysd query masterchef params

# Add external incentive
elysd tx masterchef add-external-incentive uinc 1 80 1080 1000000 --from=treasury --yes --gas=1000000
elysd query masterchef external-incentive 0

# Query pending rewards
elysd query masterchef user-pending-reward $TREASURY

# Claim rewards
elysd tx masterchef claim-rewards --pool-ids=1 --from=treasury --yes --gas=1000000

# Test flow
# - oracle setup without lifetime
# - amm pools creation
# - create external incentive on materchef (2)
# - try querying rewards
# - try claiming rewards and check state change after that
# - try swap operation and check if swap fee's correctly distributed