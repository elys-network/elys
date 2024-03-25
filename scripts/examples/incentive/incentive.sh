#!/usr/bin/env bash

# Local test
elysd tx gov submit-proposal proposal.json --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000
elysd tx gov vote 1 Yes --from=treasury --keyring-backend=test --chain-id=elystestnet-1 --yes --gas=1000000

# Query pool daily rewards
elysd query incentive pool-rewards ""
pools:
- pool_id: "1"
  reward_coins: []
  rewardsUsd: "0.000000000000000000"
- pool_id: "2"
  reward_coins:
  - amount: "162942450"
    denom: ueden
  rewardsUsd: "1868.487640182276004650"
- pool_id: "3"
  reward_coins:
  - amount: "78653292"
    denom: ueden
  rewardsUsd: "901.930122946153613244"