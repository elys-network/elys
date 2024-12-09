#!/usr/bin/env bash

# start old chain
cd /<path-to-old>/elys
elysd start

# stake Elys token
elysd query staking validators
elysd keys add user1 --keyring-backend=test --recover
# beef model palm pepper claim taste climb primary category gallery lava mimic future festival sign milk present toss basket reflect dignity frame scan hat
elysd tx bank send treasury elys1r8nfpk0t4g6r7542g99s8ymxwzpsg57alx7gue 1000000000000uelys --from=treasury --yes 
elysd tx staking delegate elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs 100000000000uelys --from=user1 --yes --gas=1000000

# claim Eden reward and commit
elysd keys show -a treasury --keyring-backend=test
elysd query commitment show-commitments elys1r8nfpk0t4g6r7542g99s8ymxwzpsg57alx7gue
elysd tx incentive withdraw-rewards --from=user1 --yes 
elysd tx commitment commit-claimed-rewards 342464 ueden --from=user1 --yes --gas=1000000

# create vesting from Eden to Elys
elysd tx commitment vest 1284240 ueden --from=user1 --yes --gas=1000000
elysd query commitment show-commitments elys1r8nfpk0t4g6r7542g99s8ymxwzpsg57alx7gue

# create liquidity pool
elysd tx amm create-pool 10uatom,10uusdc 100000000000uatom,10000000000uusdc --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false --from=treasury --yes --gas=1000000
elysd tx stablestake bond 19800000000 --from=treasury --yes --gas=1000000
elysd query commitment show-commitments elys12tzylat4udvjj56uuhu3vj2n4vgp7cf9fwna9w

# Migration proposal to "masterchef" version
elysd tx gov submit-legacy-proposal software-upgrade "masterchef" --upgrade-height=150 --deposit=10000000uelys --title="XXX" --description="XXX" --no-validate --from=treasury --fees=100000uelys --gas=200000 -y
elysd tx gov vote 1 yes --from=treasury --fees=100000uelys --gas=200000 -y
elysd tx gov vote 1 yes --from=user1 --fees=100000uelys --gas=200000 -y
elysd query gov tally 1
elysd query gov proposals

# Migrate and check if distribution module and Eden/EdenB commit work as expected
cd /<path-to-new>/elys
go install ./cmd/elysd/
elysd start

# Test functionalities after upgrade
VALIDATOR=elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs
EDEN_VAL=elysvaloper1gnmpr8vvslp3shcq6e922xr0uq4aa2w5gdzht0
EDENB_VAL=elysvaloper1wajd6ekh9u37hyghyw4mme59qmjllzuyaceanm

# estaking functionalities
elysd tx commitment commit-claimed-rewards 173430717 ueden --from=treasury --yes --gas=1000000
elysd query distribution rewards $TREASURY $EDEN_VAL
elysd query commitment show-commitments $TREASURY
elysd query staking delegations $TREASURY

elysd tx distribution withdraw-rewards $VALIDATOR --from=treasury --yes --gas=1000000
elysd tx distribution withdraw-rewards $EDEN_VAL --from=treasury --yes --gas=1000000
elysd tx distribution withdraw-rewards $EDENB_VAL --from=treasury --yes --gas=1000000

elysd tx staking delegate elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs 100000000000uelys --from=treasury --yes --gas=1000000
elysd query distribution rewards elys1r8nfpk0t4g6r7542g99s8ymxwzpsg57alx7gue $EDEN_VAL
elysd query estaking params

# masterchef functionalities
elysd query masterchef user-pending-reward elys1r8nfpk0t4g6r7542g99s8ymxwzpsg57alx7gue
elysd tx masterchef claim-rewards --pool-ids=1 --from=treasury --yes --gas=1000000
elysd tx masterchef claim-rewards --pool-ids=32767 --from=treasury --yes --gas=1000000

