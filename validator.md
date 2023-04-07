# Validator Guide

The Elys blockchain network is built on top of the Terdermint consensus engine. Terdermint leverages a group of validators tasked with appending new blocks to the blockchain. These validators engage in the consensus protocol by propagating vote messages consisting of cryptographic signatures, which are individually signed by each validator's private key.

Validators have the option to bond their staking tokens and accept delegated tokens from Elys Network token holders. The Elys Network's native token is called ELYS. Initially, Elys launched with four validators, whose selection was based on the amount of stake delegated to them. The validator candidates with the highest delegated stake are deemed the top validators and are included in the active Elys validator set.

The Elys Network protocol enables validators and their delegators to earn EDEN as block rewards and tokens as transaction fees. While transaction fees are currently payable in ELYS, any token within the Cosmos ecosystem can be accepted as a fee tender provided that it has been whitelisted by the governance system. It is important to note that validators have the flexibility to set commissions on the fees earned by their delegators, offering an additional incentive.

## Limitations

Validators who double sign, exhibit frequent offline behavior, or refuse to engage in governance are at risk of having their staked ELYS, including the ELYS delegated to them by their users, slashed. The severity of the penalty is determined based on the nature and magnitude of the offense.

## Hardware Setup

Validators must establish a physically secured operation with limited access. Co-locating in secure data centers is an ideal starting point for this purpose.

Additionally, validators should prepare to equip their datacenter with redundant power, connectivity, and storage backups. To achieve this, they should install several redundant networking boxes for fiber, firewall, and switching, as well as small servers with redundant hard drives and failover capabilities. Initially, the hardware requirements can be met with low-end datacenter gear.

It is important to note that network requirements are expected to remain low in the early stages of operation. However, as the network grows, bandwidth, CPU, and memory requirements are likely to increase. For this reason, it is recommended that validators opt for large hard drives to accommodate the storage of years of blockchain history.

## Supported Operating Systems

We provide official support for macOS and Linux operating systems exclusively in the following architectures:

- darwin/arm64
- darwin/x86_64
- linux/arm64
- linux/amd64

## Minimum System Requirements

To operate mainnet or testnet validator nodes, a machine that meets the following minimum hardware specifications is required:

- At least 4 physical CPU cores
- At least 500GB of SSD disk storage
- At least 32GB of memory (RAM)
- At least 100mbps network bandwidth

It is important to note that the server requirements may increase as the usage of the blockchain grows. As a result, you should have a plan for upgrading your server accordingly.

## Create a Dedicated Validator Website and Social Profile

To ensure transparency and provide users with information about the entity to which they are delegating their ELYS, it is recommended that validators create a dedicated website and social profile (such as on Twitter). Additionally, validators should signal their intent to become a validator on Discord. This step is essential since users will be interested in learning more about the validator before staking their ELYS.

## Engage and Seek Advice from the Validator Community

It is crucial for validators to engage in comprehensive discussions and seek advice from the wider validator community. Validators can achieve this by actively participating in discussions on Discord. This enables validators to gain insights into the finer details of validator operations and seek guidance on various aspects of the process.

## Operate a Validator Node

Gain knowledge on operating a validator node by learning the necessary skills and techniques required to run the node.

### Preparation for Using a Key Management System (KMS)

Before implementing a Key Management System (KMS), it is recommended that you follow these initial steps: Understanding the Use of a KMS.

### Validator Creation

Validators can be created by staking ELYS tokens using the node consensus public key (elysvalconspub...). The validator pubkey can be obtained by executing the following command:

```bash
elysd tendermint show-validator
```

### Validator Creation on Testnet

To create a validator on the testnet, simply execute the following command:

```bash
elysd tx staking create-validator \
  --amount=1000000uelys \
  --pubkey=$(elysd tendermint show-validator) \
  --moniker="choose a moniker" \
  --chain-id=<chain_id> \
  --commission-rate="0.05" \
  --commission-max-rate="0.10" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --gas="auto" \
  --gas-prices="0.025uelys" \
  --from=<key_name>
```

### Commission Parameters

When specifying commission parameters, the `commission-max-change-rate` is utilized to calculate the percentage point change over the commission rate. For instance, going from 1% to 2% represents a 100% rate increase, but only a 1 percentage point increase.

### Definition of Min-self-delegation

`min-self-delegation` refers to a positive integer value that denotes the minimum level of self-delegated voting power that must always be maintained by the validator. For instance, a min-self-delegation of 1000000 signifies that the validator will never have a self-delegation value lower than 1 ELYS.

### Validator Set Confirmation

Third-party explorers can be used to verify if a node has been included in the validator set.

## Editing Validator Description

Validators can edit their public description to enable delegators to identify their validator and determine which validators to delegate their stake to. It is essential to provide input for every flag below when executing the command. If a flag is excluded, the field will default to empty (`--moniker` defaults to the machine name) if the field has never been set, or remain the same if it has been set in the past.

The `<key_name>` specifies the validator being edited. If certain flags are not included, remember that the `--from` flag must be included to identify the validator to update.

The `--identity` flag can be utilized to verify the identity of validators with systems like Keybase or UPort. When using Keybase, the `--identity` flag should be populated with a 16-digit string generated from a keybase.io account. This is a secure method of verifying the validator's identity across multiple online networks. The Keybase API can retrieve the Keybase avatar, allowing validators to add a logo to their validator profile.

```bash
elysd tx staking edit-validator
--moniker="choose a moniker" \
 --website="https://elys.network" \
 --identity=9F0D65E29B4CBC8E \
 --details="To infinity and beyond!" \
 --chain-id=<chain_id> \
 --gas="auto" \
 --gas-prices="0.025uelys" \
 --from=<key_name> \
 --commission-rate="0.10"
```

### Commission-Rate Value Constraints

Note that the `commission-rate` value must satisfy the following invariants:

- It must fall between 0 and the validator's `commission-max-rate`.
- It must not exceed the validator's `commission-max-change-rate`, which denotes the maximum percentage point change rate allowed per day. In other words, a validator can only change its commission once per day and within the `commission-max-change-rate` limits.

## Validator Description Viewing

Use the following command to view the validator's information:

```bash
elysd query staking validator <account_cosmos>
```

## Validator Signing Information Tracking

To track a validator's past signatures, use the `signing-info` command as shown below:

```bash
elysd query slashing signing-info <validator-pubkey> --chain-id=<chain_id>
```

## Unjailing a Validator

If a validator is "jailed" for downtime, an `unjail` transaction must be submitted from the operator account to enable the receipt of block proposer rewards again (dependent on the zone fee distribution). Execute the command as follows:

```bash
elysd tx slashing unjail --from=<key_name> --chain-id=<chain_id>
```

## Validator Confirmation

Your validator is deemed active if the following command returns a result:

```bash
elysd query tendermint-validator-set | grep "$(elysd tendermint show-address)"
```

You should be able to locate your validator in one of Elys explorers by looking for the bech32 encoded address in the `~/.elysd/config/priv_validator.json` file.

Note that to be included in the validator set, the total voting power must exceed that of the 100th validator.

## Graceful Halting of Validator

When undertaking routine maintenance or preparing for an upcoming coordinated upgrade, it may be beneficial to systematically and gracefully halt the validator. This can be accomplished in two ways:

Setting the `halt-height` to the desired height at which the node will shut down.
Passing the `--halt-height` flag to `elysd`.

The node will gracefully shut down at the specified height after committing the block, exiting with a zero exit code.

## Common Problems

### Problem #1: My validator has voting_power: 0

If your validator has a voting power of 0, it has been jailed, which implies that it has been removed from the active validator set. Validators can be jailed for either failing to vote on 500 of the last 10000 blocks or double signing.

If the validator has been jailed for downtime, it is possible to recover the voting power. Start `elysd` again if it is not already running:

```bash
elysd start
```

Allow your full node to synchronize with the latest block. You can then unjail your validator.

Finally, check the validator to verify if the voting power has been restored:

```bash
elysd status
```

If the voting power of the validator appears to be lower than before, it implies that the validator has been slashed for downtime.

### Problem #2: Node Crashes Due to Excessive Open Files

The default maximum number of files that can be opened by a Linux process is 1024, which may be insufficient for `elysd`, resulting in process crashes. To resolve this, execute ulimit -n 4096 to increase the number of open files allowed, and restart the process with `elysd` start. If you are using a process manager like systemd to launch `elysd`, it may require additional configuration. A sample systemd file to resolve this issue is provided below:

```toml
# /etc/systemd/system/elysd.service
[Unit]
Description=Elys Node
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/go/bin/elysd start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```
