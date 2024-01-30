# Architecture Guide

This section contains documentation on the architecture of the Elys chain, including the current design and components of the system.

<details>
<summary>Click to expand/collapse</summary>

## Boilerplate Generation

The boilerplate was generated using `ignite CLI`, which provides a convenient way to generate new chains, modules, messages, and more. The initial modules that are part of the repository include `AssetProfile` and few others, both of which were generated using the `ignite CLI`.

`AssetProfile` requires all changes to go through governance proposals (i.e., adding, updating, or deleting an asset profile entry). Similarly, any modules that expose parameters must require governance proposals to update the module parameters.

## Configuration File

The repository also includes a `config.yml` file, which provides a convenient way to initiate the genesis account, set up a faucet for testnet, define initial validators, and override initial genesis states. Although `ignite` provides the network layer that allows for easy onboarding of new validators to a chain network, the `config.yml` file can be used to specify additional configurations.

In the current `config.yml` file, additional denom metadata has been defined to allow for easy setting of the ELYS amount using any exponent (decimal precision) following the EVMOS good practices. The governance params have also been overridden to reduce the voting period to 20 seconds for local test purposes. Multiple `config.yml` files can be created for each environment (local, testnet, mainnet) with their specific parameters.

## Asset Profile

### Add Entry using Gov Proposal

A proposal can be submitted to add one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "add new entries",
  "description": "add new entries",
  "messages": [
    {
      "@type": "/elys.assetprofile.MsgCreateEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2",
      "decimals": "18",
      "denom": "mytoken",
      "path": "",
      "ibcChannelId": "1",
      "ibcCounterpartyChannelId": "1",
      "displayName": "mytoken",
      "displaySymbol": "mytoken",
      "network": "",
      "address": "",
      "externalSymbol": "mytoken",
      "transferLimit": "",
      "permissions": [],
      "unitDenom": "mytoken",
      "ibcCounterpartyDenom": "mytoken",
      "ibcCounterpartyChainId": "test"
    },
    {
      "@type": "/elys.assetprofile.MsgCreateEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken3",
      "decimals": "18",
      "denom": "mytoken",
      "path": "",
      "ibcChannelId": "1",
      "ibcCounterpartyChannelId": "1",
      "displayName": "mytoken",
      "displaySymbol": "mytoken",
      "network": "",
      "address": "",
      "externalSymbol": "mytoken",
      "transferLimit": "",
      "permissions": [],
      "unitDenom": "mytoken",
      "ibcCounterpartyDenom": "mytoken",
      "ibcCounterpartyChainId": "test"
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

### Update Entry using Gov Proposal

A proposal can be submitted to update one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "update existing entries",
  "description": "update existing entries",
  "messages": [
    {
      "@type": "/elys.assetprofile.MsgUpdateEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2",
      "decimals": "18",
      "denom": "mytoken2",
      "path": "",
      "ibcChannelId": "1",
      "ibcCounterpartyChannelId": "1",
      "displayName": "mytoken2",
      "displaySymbol": "mytoken2",
      "network": "",
      "address": "",
      "externalSymbol": "mytoken2",
      "transferLimit": "",
      "permissions": [],
      "unitDenom": "mytoken2",
      "ibcCounterpartyDenom": "mytoken2",
      "ibcCounterpartyChainId": "test"
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

### Delete Entry using Gov Proposal

A proposal can be submitted to delete one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "delete entries",
  "description": "delete entries",
  "messages": [
    {
      "@type": "/elys.assetprofile.MsgDeleteEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2"
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

### CLI to Query List of Entries

To query the list of entries in the asset profile module, use the following command:

```
elysd q assetprofile list-entry
```

## Tokenomics

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

## Denom Units

The `denom_units` property is an array of objects defined in the [config.yml](./config.yml) file, with each object defining a single denomination unit. Each unit object has three properties - `denom`, `exponent`, and `aliases`.

For the ELYS token, there are three denomination units defined with aliases:

- `uelys`: This is the base unit of the ELYS token, and has no aliases.

- `melys`: This unit has an exponent of 3, which means that 1 `melys` is equal to 1000 `uelys`. It has one alias - `millielys`.

- `elys`: This unit has an exponent of 6, which means that 1 `elys` is equal to 1,000,000 `uelys`. It has no aliases.

The aliases for the `melys` unit are specified as `millielys`, which is a common prefix used to denote a thousandth of a unit. These aliases can be used interchangeably with the primary unit names in order to make the values more readable and easier to work with.

</details>

## TestNet Parameters

Here are the definitions and current values of each individual parameter of the Elys TestNet Network as of May 8th, 2023.

<details>
<summary>Click to expand/collapse</summary>

## Minting

Defines the rules for automated minting of new tokens. In the current implementation, minting is entirely disabled.

## Staking

Defines the rules for staking and delegating tokens in the network. Validators and delegators must lock their tokens for a certain period to participate in consensus and receive rewards. The `unbonding_time` parameter specifies the duration for which a validator's tokens are locked after they unbond.

- `Max_entries`: The maximum number of entries in the validator set. Current value: 7.
- `Historical_entries`: The number of entries to keep in the historical validator set. Current value: 10,000.
- `Unbonding_time`: The time period for which a validator's tokens are locked after they unbond. Current value: 1,209,600 seconds (equals to 14 days).
- `Max_validators`: The maximum number of validators that can be active at once. Current value: 100.
- `Bond_denom: The denomination used for staking tokens. Current value: `uelys`.

## Governance

Defines the rules for proposing and voting on changes to the network. To make a proposal, a minimum deposit of ELYS is required. The proposal must then go through a voting process where a certain percentage of bonded tokens must vote, and a certain percentage of those votes must be in favor of the proposal for it to pass.

- `Min_deposit`: The minimum amount of ELYS required for a proposal to enter voting. Current value: 10 ELYS.
- `Max_deposit_period`: The maximum period for which deposits can be made for a proposal. Current value: 60.
- `Quorum: The minimum percentage of total bonded tokens that must vote for a proposal to be considered valid. Current value: 33.4%.
- `Threshold`: The minimum percentage of yes votes required for a proposal to pass. Current value: 50%.
- `Veto_threshold`: The percentage of no votes required to veto a proposal. Current value: 33.4%.
- `Voting_period`: The period for which voting on a proposal is open. Current value: 60.

## Distribution

Defines the distribution of rewards and fees in the network. Block proposers receive a portion of the block rewards as an incentive to maintain the network.

- `Base_proposer_reward`: The base percentage of block rewards given to proposers. Current value: 1%.
- `Bonus_proposer_reward`: The additional percentage of block rewards given to proposers if they include all valid transactions. Current value: 4%.

## Slashing

Defines the penalties for validators who violate the network rules or fail to perform their duties. Validators who sign blocks incorrectly or go offline for too long will be penalized with a percentage of their bonded tokens being slashed. The `signed_blocks_window` parameter specifies the number of blocks used to determine a validator's uptime percentage, and the `min_signed_per_window` parameter specifies the minimum percentage of blocks that a validator must sign in each window to avoid being slashed. The `downtime_jail_duration` parameter specifies the duration for which a validator is jailed if they miss too many blocks.

- `Signed_blocks_window`: The number of blocks used to determine a validator's uptime percentage. Current value: 30,000.
- `Min_signed_per_window`: The minimum percentage of blocks that a validator must sign in each window to avoid being slashed. Current value: 5%.
- `Downtime_jail_duration`: The duration for which a validator is jailed if they miss too many blocks. Current value: 600 seconds.
- `Slash_fraction_double_sign`: The percentage of a validator's bonded tokens that are slashed if they double sign. Current value: 0.01%.
- `Slash_fraction_downtime`: The percentage of a validator's bonded tokens that are slashed if they are offline for too long. Current value: 5%.

</details>
