<!--
order: 3
-->

# Usage

## Commands

### Querying Parameters

```bash
elysd query stablestake params
```

### Querying Borrow Ratio

```bash
elysd query stablestake borrow-ratio
```

### Bonding Tokens

```bash
elysd tx stablestake bond 1000000000000uusdc --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```

### Unbonding Tokens

```bash
elysd tx stablestake unbond 500000000000uusdc --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
```
