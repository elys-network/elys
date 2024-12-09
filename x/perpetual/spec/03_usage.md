<!--
order: 3
-->

# Usage

## Transaction Commands

### Open a Position

#### Infinite Profitability

To open a long position with infinite profitability:

```shell
elysd tx perpetual open long 5 uatom 100000000uusdc --from=treasury --yes --gas=1000000
```

#### Finite Profitability

To open a short position with a specified take profit level:

```shell
elysd tx perpetual open short 5 uatom 100000000uusdc --take-profit 100 --from=treasury --yes --gas=1000000
```

### Close a Position

To close an existing position:

```shell
elysd tx perpetual close 1 10000000 --from=treasury --yes --gas=1000000
```

### Manage Whitelist

#### Whitelist an Address

To whitelist an address for trading in the perpetual module:

```shell
elysd tx perpetual whitelist ADDRESS --from=treasury --yes --gas=1000000
```

#### Dewhitelist an Address

To remove an address from the whitelist:

```shell
elysd tx perpetual dewhitelist ADDRESS --from=treasury --yes --gas=1000000
```

### Update Module Parameters

To update the parameters of the perpetual module:

```shell
elysd tx perpetual update-params [OPTIONS] --from=treasury --yes --gas=1000000
```

## Query Commands

### Query Module Parameters

To query the current parameters of the perpetual module:

```shell
elysd query perpetual params
```

### Query All Positions

To query all positions:

```shell
elysd query perpetual get-positions
```

### Query Positions by Pool

To query positions for a specific pool:

```shell
elysd query perpetual get-positions-by-pool [amm_pool_id]
```

### Query Positions for an Address

To query positions for a specific address:

```shell
elysd query perpetual get-positions-for-address [address]
```

### Query Module Status

To query the current status of the perpetual module:

```shell
elysd query perpetual get-status
```

### Query Whitelisted Addresses

To query all whitelisted addresses:

```shell
elysd query perpetual get-whitelist
```

### Check if Address is Whitelisted

To check if a specific address is whitelisted:

```shell
elysd query perpetual is-whitelisted [address]
```

### Query Pool List

To query all pools:

```shell
elysd query perpetual list-pool
```

### Query Specific Pool

To query details of a specific pool:

```shell
elysd query perpetual show-pool [index]
```

### Query MTP

To query a specific MTP (margin trading position):

```shell
elysd query perpetual get-mtp [address] [id]
```

### Open Estimation

To query an open estimation:

```shell
elysd query perpetual open-estimation [position] [leverage] [trading-asset] [collateral]
```

## Examples

### Open a Long Position

```shell
elysd tx perpetual open long 5 uatom 100000000uusdc --from=treasury --yes --gas=1000000
```

### Open a Short Position with Take Profit

```shell
elysd tx perpetual open short 5 uatom 100000000uusdc --take-profit 100 --from=treasury --yes --gas=1000000
```

### Close a Position

```shell
elysd tx perpetual close 1 10000000 --from=treasury --yes --gas=1000000
```

### Whitelist an Address

```shell
elysd tx perpetual whitelist elys1qv4k64xr6nhcxgnzq0l8t8wy9s9d4s9d6w93lz --from=treasury --yes --gas=1000000
```

### Dewhitelist an Address

```shell
elysd tx perpetual dewhitelist elys1qv4k64xr6nhcxgnzq0l8t8wy9s9d4s9d6w93lz --from=treasury --yes --gas=1000000
```

### Update Parameters

```shell
elysd tx perpetual update-params --maxLeverage=10 --maintenanceMarginRatio=0.05 --fundingRateInterval=60 --from=treasury --yes --gas=1000000
```

### Query a Specific Position

```shell
elysd query perpetual get-mtp elys1qv4k64xr6nhcxgnzq0l8t8wy9s9d4s9d6w93lz 1
```

### Query All Positions

```shell
elysd query perpetual get-positions
```

### Query Positions by Pool

```shell
elysd query perpetual get-positions-by-pool 1
```

### Query Positions for an Address

```shell
elysd query perpetual get-positions-for-address elys1qv4k64xr6nhcxgnzq0l8t8wy9s9d4s9d6w93lz
```

### Query Module Status

```shell
elysd query perpetual get-status
```

### Query Whitelisted Addresses

```shell
elysd query perpetual get-whitelist
```

### Check if Address is Whitelisted

```shell
elysd query perpetual is-whitelisted elys1qv4k64xr6nhcxgnzq0l8t8wy9s9d4s9d6w93lz
```

### Query All Pools

```shell
elysd query perpetual list-pool
```

### Query Specific Pool

```shell
elysd query perpetual show-pool 1
```

### Open Estimation

```shell
elysd query perpetual open-estimation long 5 uatom 100000000uusdc
```
