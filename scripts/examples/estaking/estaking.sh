#!/usr/bin/env bash

TREASURY=$(elysd keys show -a treasury --keyring-backend=test)

elysd query commitment show-commitments $TREASURY
elysd query bank balances $TREASURY

elysd query estaking params

elysd query staking validators
VALIDATOR=elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs
EDEN_VAL=elysvaloper1gnmpr8vvslp3shcq6e922xr0uq4aa2w5gdzht0
EDENB_VAL=elysvaloper1wajd6ekh9u37hyghyw4mme59qmjllzuyaceanm

elysd tx staking delegate $VALIDATOR 1000000000000uelys --from=treasury --yes --gas=1000000
elysd query distribution rewards $TREASURY $VALIDATOR
elysd query distribution rewards $TREASURY $EDEN_VAL
elysd query distribution rewards $TREASURY $EDENB_VAL

elysd tx distribution withdraw-rewards $VALIDATOR --from=treasury --yes --gas=1000000
elysd tx distribution withdraw-rewards $EDEN_VAL --from=treasury --yes --gas=1000000
elysd tx distribution withdraw-rewards $EDENB_VAL --from=treasury --yes --gas=1000000

elysd tx commitment commit-claimed-rewards 503544 ueden --from=treasury --yes --gas=1000000
elysd tx commitment commit-claimed-rewards 1678547 uedenb --from=treasury --yes --gas=1000000

# Pay 10000 uusdc as gas fees
elysd tx staking delegate $VALIDATOR 1000uelys --fees=10000uusdc --from=treasury --yes --gas=1000000

elysd query estaking rewards $TREASURY
elysd tx estaking withdraw-all-rewards --from=validator
