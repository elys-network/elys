# elys

**elys** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Ignite CLI has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release

To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install

To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/elys-network/elys@latest! | sudo bash
```

`elys-network/elys` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)

## Boilerplate Generation

The boilerplate was generated using `ignite CLI`, which provides a convenient way to generate new chains, modules, messages, and more. The initial modules that are part of the repository include `AssetProfile` and `LiquidityProvider`, both of which were generated using the `ignite CLI`.

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
      "@type": "/elysnetwork.elys.assetprofile.MsgCreateEntry",
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
      "@type": "/elysnetwork.elys.assetprofile.MsgCreateEntry",
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
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### Update Entry using Gov Proposal

A proposal can be submitted to update one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "update existing entries",
  "description": "update existing entries",
  "messages": [
    {
      "@type": "/elysnetwork.elys.assetprofile.MsgUpdateEntry",
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
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### Delete Entry using Gov Proposal

A proposal can be submitted to delete one or multiple entries in the asset profile module. The proposal must be in the following format:

```json
{
  "title": "delete entries",
  "description": "delete entries",
  "messages": [
    {
      "@type": "/elysnetwork.elys.assetprofile.MsgDeleteEntry",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "baseDenom": "mytoken2"
    }
  ],
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
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
      "@type": "/elysnetwork.elys.tokenomics.MsgUpdateGenesisInflation",
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
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### CLI to Query the Genesis Inflation parameters

To query the gensis inflation parameters in the tokenomics module, use the following command:

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
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "AtomStakers",
      "amount": "9999999"
    },
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "RowanStakersLP",
      "amount": "9999999"
    },
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "Juno",
      "amount": "9999999"
    },
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "Osmo",
      "amount": "9999999"
    },
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "Evmos",
      "amount": "9999999"
    }
  ],
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### Update Airdrop entry using Gov Proposal

A proposal can be submitted to update one or multiple airdrop entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "update existing entries",
  "description": "update existing entries",
  "messages": [
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgUpdateAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "AtomStakers",
      "amount": "9999999"
    }
  ],
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### Delete Airdrop entry using Gov Proposal

A proposal can be submitted to delete one or multiple airdrop entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "delete airdrop entries",
  "description": "delete airdrop entries",
  "messages": [
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgDeleteAirdrop",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "intent": "AtomStakers"
    }
  ],
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
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
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateTimeBasedInflation",
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
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateTimeBasedInflation",
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
      "@type": "/elysnetwork.elys.tokenomics.MsgCreateTimeBasedInflation",
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
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### Update Time-Based-Inflation entry using Gov Proposal

A proposal can be submitted to update one or multiple time-based-inflation entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "update existing time-based-inflation entries",
  "description": "update existing time-based-inflation entries",
  "messages": [
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgUpdateTimeBasedInflation",
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
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### Delete Time-Based-Inflation entry using Gov Proposal

A proposal can be submitted to delete one or multiple time-based-inflation entries in the tokenomics module. The proposal must be in the following format:

```json
{
  "title": "delete time-based-inflation entries",
  "description": "delete time-based-inflation entries",
  "messages": [
    {
      "@type": "/elysnetwork.elys.tokenomics.MsgDeleteTimeBasedInflation",
      "authority": "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
      "startBlockHeight": "12614402",
      "endBlockHeight": "18921602"
    }
  ],
  "deposit": "10000000000000000000aelys"
}
```

To submit a proposal, use the following command:

```
elysd tx gov submit-proposal /tmp/proposal.json --from alice --yes
```

To vote on a proposal, use the following command:

```
elysd tx gov vote 1 yes --from alice --yes
```

### CLI to Query List of Time-Based-Inflation entries

To query the list of the time-based-inflation entries in the tokenomics module, use the following command:

```
elysd q tokenomics list-time-based-inflation
```
