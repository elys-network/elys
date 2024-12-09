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

### Querying Parameters

```bash
elysd query commitment params
```

### Committing Tokens

```bash
elysd tx commitment commit-claimed-rewards 1000000000uelys --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Uncommitting Tokens

```bash
elysd tx commitment uncommit-tokens 500000000uelys --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Vesting Tokens

```bash
elysd tx commitment vest 100000000uelys --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Claiming Vested Tokens

```bash
elysd tx commitment claim-vesting --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Staking Tokens

```bash
VALIDATOR=elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs
elysd tx commitment stake 100000000uelys $VALIDATOR --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Unstaking Tokens

```bash
VALIDATOR=elysvaloper12tzylat4udvjj56uuhu3vj2n4vgp7cf9pwcqcs
elysd tx commitment unstake 50000000uelys $VALIDATOR --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```
