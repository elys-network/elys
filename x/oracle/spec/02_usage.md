<!--
order: 2
-->

# Usage

## Commands

### Querying Parameters and Band Price Results

```bash
elysd query oracle params
elysd query oracle band-price-result [request_id]
elysd query oracle last-band-request-id
```

### Querying Asset Info

```bash
elysd query oracle list-asset-info
elysd query oracle show-asset-info [denom]
```

### Querying Prices

```bash
elysd query oracle list-price
elysd query oracle show-price [asset] --source="[source]" --timestamp=[timestamp]
elysd query oracle show-price [asset] --source="[source]"
elysd query oracle show-price [asset]
```

### Managing Price Feeders

```bash
elysd query oracle list-price-feeder
elysd query oracle show-price-feeder [feeder_address]

elysd tx oracle set-price-feeder [feeder_address] [is_active] --from=[key] --chain-id=[chain-id] --broadcast-mode=block --yes
elysd tx oracle delete-price-feeder [feeder_address] --from=[key] --chain-id=[chain-id] --broadcast-mode=block --yes
```

### Feeding Prices

```bash
elysd tx oracle feed-price [asset] [price] [source] --from=[provider] --chain-id=[chain-id] --broadcast-mode=block --yes
elysd tx oracle feed-multiple-prices [prices-json] --from=[creator] --chain-id=[chain-id] --broadcast-mode=block --yes
```

### Managing Asset Info

```bash
elysd tx oracle add-asset-info-proposal [denom] [display] [band_ticker] [elys_ticker] [decimal] --title="[title]" --description="[description]" --deposit="[deposit_amount]" --from=[authority] --chain-id=[chain-id] --broadcast-mode=block --yes
elysd tx oracle remove-asset-info-proposal [denom] --title="[title]" --description="[description]" --deposit="[deposit_amount]" --from=[authority] --chain-id=[chain-id] --broadcast-mode=block --yes
```

### Submitting Governance Proposals

```bash
elysd query gov proposals
elysd tx gov vote [proposal_id] yes --from=[voter] --chain-id=[chain-id] --broadcast-mode=block --yes
```

### Requesting Band Price

```bash
elysd tx oracle request-band-price --from=[key] --chain-id=[chain-id] --broadcast-mode=block --yes
```
