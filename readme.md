# Elys

**Elys** is a blockchain built using Cosmos SDK, Tendermint and Ignite.

| Parameter    | Value                                |
| ------------ | ------------------------------------ |
| Chain ID     | elystestnet-1                        |
| Denomination | uelys                                |
| Decimals     | 6 (1 elys= 1000000uelys)             |
| Version      | v0.4.0                               |
| RPC Endpoint | https://rpc.testnet.elys.network:443 |

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts Elys in development.

## Installation

### With Ignite (Experimental)

To install the latest version of Elys binary, execute the following command on your machine:

```
curl https://get.ignite.com/elys-network/elys@latest! | sudo bash
```

### Manual Installation With Makefile (Recommended)

This section provides a step-by-step guide on how to build the Elys Chain binary from the source code using the provided makefile. The makefile automates the build process and generates a binary executable that can be run on your local machine.

<details>
<summary>Click to expand/collapse</summary>

1. Clone the Elys chain repository:

```bash
git clone https://github.com/elys-network/elys.git
```

2. Navigate to the cloned repository:

```bash
cd elys
```

3. Optionally, checkout the specific branch or tag you want to build:

```bash
git checkout <version>
```

note: 'latest' is currently not recognized but will be supported in the next version (eg use 'git checkout v.0.2.3')

4. Ensure that you have the necessary dependencies installed. For instance, on Ubuntu you need to install the `make` tool:

```bash
sudo apt-get install --yes make
```

5. Run the `make build` command to build the binary:

```bash
make build
```

6. The binary will be generated in the `./build` directory. You can run the binary using the following command:

```bash
./build/elysd
```

You can also use the `make install` command to install the binary in the `bin` directory of your `GOPATH`.

</details>

## Validator Guide

The validator guide is accessible [here](./validator.md).

## Network Launch

This section provides a step-by-step guide on how to launch a new network, such as a testnet, for Elys. The guide includes instructions on how to use Ignite commands to set up and configure the new network.

<details>
<summary>Click to expand/collapse</summary>

### Coordinator Configuration

To publish the information about Elys chain as a coordinator, run the following command:

```
ignite network chain publish github.com/elys-network/elys --tag v0.1.0 --chain-id elystestnet-1 --account-balance 10000000000uelys
```

### Validator Configuration

This documentation presupposes the validator node is currently operational on `Ubuntu 22.04.2 LTS`.

#### Prerequisites

Before launching a validator node, a set of tools must be installed.

To install the `build-essential` package, enter the following command:

```
sudo apt install build-essential
```

Install `go` version `1.19`

```
cd /tmp
wget https://go.dev/dl/go1.19.7.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.19.7.linux-amd64.tar.gz
```

Append the following line to the end of the `~/.bashrc` file:

```
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
```

Run the following command:

```
go version
```

This should return the following output:

```
go version go1.19.7 linux/amd64
```

Install `ignite-cli`

Enter the following command to install the `ignite-cli` command:

```
curl https://get.ignite.com/cli! | bash
```

Then run the following command:

```
ignite network
```

Install the latest version of Elys binary by running the following command:

```
curl https://get.ignite.com/elys-network/elys@latest! | sudo bash
```

Enter the following command to initialize the validator node and request to join the network:

```
ignite network chain init 12
ignite network chain join 12 --amount 95000000uelys
```

The coordinator will then have to approve the validator requests with the following commands:

```
ignite network request list 12
ignite network request approve 12 <REQUEST_ID>,<REQUEST_ID>
```

Once all the validators needed for the validator set are approved, to launch the chain use the following command:

```
ignite network chain launch 12
```

Each validator is now ready to prepare their nodes for launch by using this command:

```
ignite network chain prepare 12
```

The output of this command will show a command that a validator would use to launch their node such as:

```
elysd start --home $HOME/spn/12 2> elysd.log &
```

A systemd service can be created to auto-start the `elysd` service.

Create the new file `/etc/systemd/system/elysd.service` with this content:

```
[Unit]
Description=Elysd Service
Wants=network.target
After=network.target

[Service]
Environment=HOME=/home/ubuntu
Type=simple
Restart=on-failure
WorkingDirectory=/home/ubuntu
SyslogIdentifier=elysd.user-daemon
ExecStart=/home/ubuntu/go/bin/elysd start --home spn/12 2>&1
ExecStop=/usr/bin/pkill elysd

[Install]
WantedBy=multi-user.target
```

Then you can use those commands to enable and start the service:

```
sudo systemctl enable elysd.service
sudo systemctl start elysd.service
```

You can check the status of the service at any time using this command:

```
sudo systemctl status elysd.service
```

Or follow the service logs by using this command:

```
sudo journalctl -u elysd.service -f
```

</details>

## Architecture

This section contains documentation on the architecture of the Elys chain, including the current design and components of the system.

<details>
<summary>Click to expand/collapse</summary>

### Boilerplate Generation

The boilerplate was generated using `ignite CLI`, which provides a convenient way to generate new chains, modules, messages, and more. The initial modules that are part of the repository include `AssetProfile` and `LiquidityProvider`, both of which were generated using the `ignite CLI`.

`AssetProfile` requires all changes to go through governance proposals (i.e., adding, updating, or deleting an asset profile entry). Similarly, any modules that expose parameters must require governance proposals to update the module parameters.

### Configuration File

The repository also includes a `config.yml` file, which provides a convenient way to initiate the genesis account, set up a faucet for testnet, define initial validators, and override initial genesis states. Although `ignite` provides the network layer that allows for easy onboarding of new validators to a chain network, the `config.yml` file can be used to specify additional configurations.

In the current `config.yml` file, additional denom metadata has been defined to allow for easy setting of the ELYS amount using any exponent (decimal precision) following the EVMOS good practices. The governance params have also been overridden to reduce the voting period to 20 seconds for local test purposes. Multiple `config.yml` files can be created for each environment (local, testnet, mainnet) with their specific parameters.

### Asset Profile

#### Add Entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### Update Entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### Delete Entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### CLI to Query List of Entries

To query the list of entries in the asset profile module, use the following command:

```
elysd q assetprofile list-entry
```

### Tokenomics

#### Set Genesis Inflation parameters using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### CLI to Query the Genesis Inflation parameters

To query the gensis inflation parameters in the tokenomics module, use the following command:

```
elysd q tokenomics show-genesis-inflation
```

#### Add Airdrop entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### Update Airdrop entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### Delete Airdrop entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### CLI to Query List of Airdrop entries

To query the list of airdrop entries in the tokenomics module, use the following command:

```
elysd q tokenomics list-airdrop
```

#### Add Time-Based-Inflation entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### Update Time-Based-Inflation entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### Delete Time-Based-Inflation entry using Gov Proposal

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
  "deposit": "10000000uelys"
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

#### CLI to Query List of Time-Based-Inflation entries

To query the list of the time-based-inflation entries in the tokenomics module, use the following command:

```
elysd q tokenomics list-time-based-inflation
```

</details>

## Release

To release a new version of Elys, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

## Learn more

- [Twitter](https://twitter.com/elys_network)
- [TestNet Explorer](https://testnet.elys.network)
- [Developer Chat](https://discord.gg/3JtgtGJ3By)
- [Github](https://github.com/elys-network)
