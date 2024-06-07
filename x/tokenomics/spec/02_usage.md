<!--
order: 2
-->

# Usage

## Commands

### Set Genesis Inflation parameters using Gov Proposal

A proposal can be submitted to set the genesis inflation parameters in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "set new genesis inflation params",
  "description": "set new genesis inflation params",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgUpdateGenesisInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "inflation": {
        "lmRewards": "9999999",
        "icsStakingRewards": "9999999",
        "communityFund": "9999999",
        "strategicReserve": "9999999",
        "teamTokensVested": "9999999"
      },
      "seedVesting": "9999999",
      "strategicSalesVesting": "9999999"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### CLI to Query the Genesis Inflation parameters

To query the genesis inflation parameters in the tokenomics module, use the following command:

```
elysd q tokenomics show-genesis-inflation
```

### Add Airdrop entry using Gov Proposal

A proposal can be submitted to add one or multiple airdrop entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "add new airdrop entries",
  "description": "add new airdrop entries",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "AtomStakers",
      "amount": "9999999"
    },
    {
      "@type": "/elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "RowanStakersLP",
      "amount": "9999999"
    },
    {
      "@type": "/elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "Juno",
      "amount": "9999999"
    },
    {
      "@type": "/elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "Osmo",
      "amount": "9999999"
    },
    {
      "@type": "/elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "Evmos",
      "amount": "9999999"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### Update Airdrop entry using Gov Proposal

A proposal can be submitted to update one or multiple airdrop entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "update existing entries",
  "description": "update existing entries",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgUpdateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "AtomStakers",
      "amount": "9999999"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### Delete Airdrop entry using Gov Proposal

A proposal can be submitted to delete one or multiple airdrop entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "delete airdrop entries",
  "description": "delete airdrop entries",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgDeleteAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "AtomStakers"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### CLI to Query List of Airdrop entries

To query the list of airdrop entries in the tokenomics module, use the following command:

```
elysd q tokenomics list-airdrop
```

### Add Time-Based-Inflation entry using Gov Proposal

A proposal can be submitted to add one or multiple time-based-inflation entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "add new time-based-inflation entries",
  "description": "add new time-based-inflation entries",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgCreateTimeBasedInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "startBlockHeight": "1",
      "endBlockHeight": "6307200",
      "description": "1st Year Inflation",
      "inflation": {
        "lmRewards": "9999999",
        "icsStakingRewards": "9999999",
        "communityFund": "9999999",
        "strategicReserve": "9999999",
        "teamTokensVested": "9999999"
      }
    },
    {
      "@type": "/elys.tokenomics.MsgCreateTimeBasedInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "startBlockHeight": "6307201",
      "endBlockHeight": "6307200",
      "description": "2nd Year Inflation",
      "inflation": {
        "lmRewards": "9999999",
        "icsStakingRewards": "9999999",
        "communityFund": "9999999",
        "strategicReserve": "9999999",
        "teamTokensVested": "9999999"
      }
    },
    {
      "@type": "/elys.tokenomics.MsgCreateTimeBasedInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "startBlockHeight": "12614402",
      "endBlockHeight": "18921602",
      "description": "3rd Year Inflation",
      "inflation": {
        "lmRewards": "9999999",
        "icsStakingRewards": "9999999",
        "communityFund": "9999999",
        "strategicReserve": "9999999",
        "teamTokensVested": "9999999"
      }
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### Update Time-Based-Inflation entry using Gov Proposal

A proposal can be submitted to update one or multiple time-based-inflation entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "update existing time-based-inflation entries",
  "description": "update existing time-based-inflation entries",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgUpdateTimeBasedInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "startBlockHeight": "12614402",
      "endBlockHeight": "18921602",
      "description": "Updated 3rd Year Inflation",
      "inflation": {
        "lmRewards": "9999999",
        "icsStakingRewards": "9999999",
        "communityFund": "9999999",
        "strategicReserve": "9999999",
        "teamTokensVested": "9999999"
      }
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### Delete Time-Based-Inflation entry using Gov Proposal

A proposal can be submitted to delete one or multiple time-based-inflation entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "delete time-based-inflation entries",
  "description": "delete time-based-inflation entries",
  "messages": [
    {
      "@type": "/elys.tokenomics.MsgDeleteTimeBasedInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "startBlockHeight": "12614402",
      "endBlockHeight": "18921602"
    }
  ],
  "deposit": "10000000uelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from walletname --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from walletname --yes
```

### CLI to Query List of Time-Based-Inflation entries

To query the list of the time-based-inflation entries in the tokenomics module, use the following command:

```
elysd q tokenomics list-time-based-inflation
```
