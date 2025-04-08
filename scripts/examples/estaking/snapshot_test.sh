#!/usr/bin/env bash

elysd query estaking params

elysd query staking validators
VALIDATOR=elysvaloper1dd34v384hdfqgajkg0jzp0y5k6qlvhltr73as2
EDEN_VAL=elysvaloper1gnmpr8vvslp3shcq6e922xr0uq4aa2w5gdzht0
EDENB_VAL=elysvaloper1wajd6ekh9u37hyghyw4mme59qmjllzuyaceanm

elysd tx staking delegate $VALIDATOR 10000000000000uelys --from=validator --yes --gas=1000000
elysd query distribution rewards elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg $VALIDATOR
elysd query distribution rewards elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg $EDEN_VAL
elysd query distribution rewards elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg $EDENB_VAL

elysd tx estaking withdraw-all-rewards --from=validator --yes --gas=1000000

elysd tx commitment commit-claimed-rewards 158762097 ueden --from=validator --yes --gas=1000000
elysd tx commitment commit-claimed-rewards 1066283235 uedenb --from=validator --yes --gas=1000000
elysd tx commitment uncommit-tokens 35235693 ueden --from=validator --yes --gas=1000000
elysd tx commitment commit-claimed-rewards 35235693 ueden --from=validator --yes --gas=1000000
elysd tx commitment uncommit-tokens 304152385 uedenb --from=validator --yes --gas=1000000
elysd tx commitment uncommit-tokens 70471386 ueden --from=validator --yes --gas=1000000

elysd query commitment show-commitments elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg
elysd query estaking rewards elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg
elysd tx estaking withdraw-all-rewards --from=validator

elysd query incentive all-program-rewards elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg
