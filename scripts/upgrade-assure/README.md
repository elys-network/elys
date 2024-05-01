# Localnet Upgrade and Data Migration Guide for ELYS Network

This document offers a detailed framework for upgrading the localnet used in the ELYS network, encompassing version changes, data migrations, and utilizing testnet snapshots.

## Prerequisites

Ensure the following prerequisites are met before proceeding with the upgrade and migration:

- Git is installed on your machine.
- Go programming language is installed.
- Access to the project repository.
- `curl` and `make` utilities are installed.
- Appropriate permissions to execute the commands listed below.

## Upgrade Steps

### Step 1: Checkout the Branch containing the new implementation

Switch to the branch containing the new implementation:

```bash
git checkout <branch_name>
```

### Step 2: Create a New Tag using semantic versioning and Install It.

i.e for tag v0.31.0, tag the new release and install it:

```bash
git tag v0.31.0
make install
git checkout v0.31.0
```

### Step 3: Retrieve the current binary depending on your OS

i.e In case of current Binary is v0.30.0:
For MacOS/Darwin users:

```bash
curl -L https://github.com/elys-network/elys/releases/download/v0.30.0/elysd-v0.30.0-darwin-arm64 -o /tmp/elysd-v0.30.0
```

For Linux users:

```bash
curl -L https://github.com/elys-network/elys/releases/download/v0.30.0/elysd-v0.30.0-linux-amd64 -o /tmp/elysd-v0.30.0
```

### Step 4: Retrieve Testnet Snapshot

Fetch the latest testnet snapshot, necessary for data migration using `curl`:

```bash
curl -L https://snapshots-testnet.stake-town.com/elys/elystestnet-1_latest.tar.lz4 -o /tmp/snapshot.tar.lz4
```

### Step 5: Checkout the latest working tag available in https://github.com/elys-network/elys/releases/tag

Switch to the previously stable version.

```bash
git checkout <tag_name>
```

### Step 6: Run Upgrade-Assure Script

#### 6a: Initial Run

Run the upgrade-assure script without starting the node:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --skip-node-start
```

Notice that /tmp/elysd-v0.30.0 is the current binary retrieved in step 3.

#### 6b: Handle Potential Errors

Address any type errors, such as difficulties in unmarshaling strings into integer fields in Go struct fields.

#### 6c: Update the Script

Modify `scripts/upgrade-assure/types.go` to reflect data structure changes necessary to resolve type errors.

The `types.go` file employs the `elys` data structure types to serialize the genesis state into JSON format for initializing localnet. This file predominantly handles conversion issues where Go struggles with fields defined as integers. To address this, such fields are overridden as `json.Number`.

During the `read-genisis-file` step of the `upgrade-assure` process, if parsing of the genesis JSON file fails, an error is returned. This issue generally arises from integer fields that must be redefined to `json.Number`.

#### 6d: Retry Upgrade-Assure

Repeat the process after updating the script:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --skip-node-start
```

Notice that /tmp/elysd-v0.30.0 is the current binary retrieved in step 3.

### Step 7: Checkout to Latest Changes Branch you used in step 1

Switch back to the main branch to incorporate the latest changes:

```bash
git checkout <branch_name>
```

### Step 8: Final Upgrade Command

Execute the final upgrade command to complete the upgrade process:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --skip-snapshot --skip-chain-init
```

Notice that /tmp/elysd-v0.30.0 is the current binary retrieved in step 3.

### Step 9: Run the chain

1. Run the first node with `elysd start`
2. After running and the log prints "No addresses to dial" you have to start the second node:
   `elysd start --home ~/.elys2 --rpc.laddr tcp://127.0.0.1:26667 --p2p.laddr tcp://0.0.0.0:26666`.

### Step 10 (optional):

You can make a backup copy of .elys and .elys2 folders available in your home directory in case you want
to start a fresh copy of the chain without going to this process again.

### Step 11: Start the Nodes manually (optional)

If something went wrong while you were starting the node at step 8, you can start the nodes manually with the new binary by using the following command:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --only-start-with-new-binary
```

Notice that /tmp/elysd-v0.30.0 is the current binary retrieved in step 3.

## Testnet Snapshots Usage

**Snapshot Sources and Installation Procedures:**

### High Stakes Testnet

- **Snapshot Source:** [Download the latest snapshot for High Stakes Testnet](https://tools.highstakes.ch/files/elys.tar.gz).
- **Installation Commands:**
  ```bash
  make install
  go run ./scripts/upgrade-assure/... https://tools.highstakes.ch/files/elys.tar.gz ~/go/bin/elysd ~/go/bin/elysd --skip-proposal
  ```

### Stake Town Testnet

- **Snapshot Source:** [Download the latest snapshot for Stake Town Testnet](https://snapshots-testnet.stake-town.com/elys/elystestnet-1_latest.tar.lz4).
- **Installation Commands:**
  ```bash
  make install
  go run ./scripts/upgrade-assure/... https://snapshots-testnet.stake-town.com/elys/elystestnet-1_latest.tar.lz4 ~/go/bin/elysd ~/go/bin/elysd --skip-proposal
  ```

## Troubleshooting

**Common Issues and Solutions:**

- **Memory Limitation:** Address processes terminated due to insufficient RAM by creating a swap file as detailed [here](https://wiki.manjaro.org/index.php?title=Swap#Using_a_Swapfile).
- **Timeout Issues:** Modify timeout settings for node responsiveness or block processing delays:

```bash
--timeout-wait-for-node=600  # Time in seconds
--timeout-next-block=15       # Time in minutes
```

**My nodes crashed after the upgrade. What should I do?**

Run the following command to start the nodes manually:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 ~/go/bin/elysd ~/go/bin/elysd --only-start-with-new-binary
```

**Debug Mode**

By default the nodes run in `info` mode. To enable debug mode, add the following flag to the command:

```bash
LOG_LEVEL=debug go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 ~/go/bin/elysd ~/go/bin/elysd --only-start-with-new-binary
```
