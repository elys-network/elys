<!--
order: 2
-->

# Usage

## Commands

### Querying Parameters

```bash
elysd query leveragelp params
```

### Querying Pools

```bash
elysd query leveragelp pools
elysd query leveragelp pool <pool_id>
```

### Querying Positions

```bash
elysd query leveragelp positions
elysd query leveragelp position <address> <id>
elysd query leveragelp positions-by-pool <pool_id>
elysd query leveragelp positions-for-address <address>
```

### Querying Status and Whitelist

```bash
elysd query leveragelp status
elysd query leveragelp whitelist
elysd query leveragelp is-whitelisted <address>
```

### Opening a Leveraged Position

```bash
elysd tx leveragelp open --from <address> --collateral-asset <asset> --collateral-amount <amount> --amm-pool-id <pool_id> --leverage <leverage> --stop-loss-price <price> --chain-id <chain_id> --yes
```

### Closing a Leveraged Position

```bash
elysd tx leveragelp close --from <address> --id <position_id> --lp-amount <amount> --chain-id <chain_id> --yes
```

### Whitelisting and Dewhitelisting Addresses

```bash
elysd tx leveragelp whitelist --from <authority_address> --whitelisted-address <address> --chain-id <chain_id> --yes
elysd tx leveragelp dewhitelist --from <authority_address> --whitelisted-address <address> --chain-id <chain_id> --yes
```
