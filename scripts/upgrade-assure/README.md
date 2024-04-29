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

### Step 1: Checkout the Implementation Branch

Switch to the branch containing the new implementation:

```bash
git checkout estaking_impl
```

### Step 2: Create a New Tag and Install

Tag the new release and install it:

```bash
git tag v0.31.0
make install
```

### Step 3: Retrieve the current binary depending on your OS

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

### Step 5: Checkout the Previous Version

Switch to the previously stable version:

```bash
git checkout v0.30.0
```

### Step 6: Run Upgrade-Assure Script

#### 6a: Initial Run

Run the upgrade-assure script without starting the node:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --skip-node-start
```

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

### Step 7: Checkout to Latest Changes Branch

Switch back to the main branch to incorporate the latest changes:

```bash
git checkout main
```

### Step 8: Final Upgrade Command

Execute the final upgrade command to complete the upgrade process:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --skip-snapshot --skip-chain-init
```

### Step 9: Start the Nodes manually (optional)

If something went wrong while you were starting the node at step 8, you can start the nodes manually with the new binary by using the following command:

```bash
go run ./scripts/upgrade-assure/... /tmp/snapshot.tar.lz4 /tmp/elysd-v0.30.0 ~/go/bin/elysd --only-start-with-new-binary
```

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
