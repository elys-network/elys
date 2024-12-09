<!--
order: 2
-->

# Usage

## Commands

### Querying Commitments and Balances

```bash
TREASURY=$(elysd keys show -a treasury --keyring-backend=test)

elysd query commitment show-commitments $TREASURY
elysd query bank balances $TREASURY
```

### Querying Staking and Validators

```bash
elysd query estaking params
elysd query staking validators
```

### Delegating Tokens

```bash
VALIDATOR=elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs
elysd tx staking delegate $VALIDATOR 1000000000000uelys --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Querying and Withdrawing Rewards

```bash
elysd query distribution rewards $TREASURY $VALIDATOR
elysd query distribution rewards $TREASURY $EDEN_VAL
elysd query distribution rewards $TREASURY $EDENB_VAL

elysd tx distribution withdraw-rewards $VALIDATOR --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
elysd tx distribution withdraw-rewards $EDEN_VAL --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
elysd tx distribution withdraw-rewards $EDENB_VAL --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Committing Claimed Rewards

```bash
elysd tx commitment commit-claimed-rewards 503544 ueden --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
elysd tx commitment commit-claimed-rewards 1678547 uedenb --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Delegating

```bash
elysd tx staking delegate $VALIDATOR 1000uelys --fees=10000uusdc --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Querying eStaking Rewards

```bash
elysd query estaking rewards $TREASURY
```

### Withdrawing All Rewards

```bash
elysd tx estaking withdraw-all-rewards --from=validator --chain-id=elysicstestnet-1
```
