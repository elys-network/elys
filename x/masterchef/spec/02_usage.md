<!--
order: 2
-->

# Usage

## Commands

### Querying Parameters

```bash
elysd query masterchef params
```

### Querying External Incentives

```bash
elysd query masterchef external-incentive [incentive_id]
```

### Querying Pool Information

```bash
elysd query masterchef pool-info [pool_id]
```

### Querying Pool Reward Information

```bash
elysd query masterchef pool-reward-info [pool_id] [reward_denom]
```

### Querying User Reward Information

```bash
elysd query masterchef user-reward-info [user_address] [pool_id] [reward_denom]
```

### Querying User Pending Rewards

```bash
elysd query masterchef user-pending-reward [user_address]
```

### Querying Stable Stake APR

```bash
elysd query masterchef stable-stake-apr [denom]
```

### Querying Pool APRs

```bash
elysd query masterchef pool-aprs
```

### Adding External Reward Denom

```bash
elysd tx masterchef add-external-reward-denom [authority] [reward_denom] [min_amount] [supported] --from=[key_name] --chain-id=[chain_id] --yes --gas=[gas_limit]
```

### Adding External Incentive

```bash
elysd tx masterchef add-external-incentive [sender] [reward_denom] [pool_id] [from_block] [to_block] [amount_per_block] --from=[key_name] --chain-id=[chain_id] --yes --gas=[gas_limit]
```

### Claiming Rewards

```bash
elysd tx masterchef claim-rewards [sender] [pool_ids] --from=[key_name] --chain-id=[chain_id] --yes --gas=[gas_limit]
```
