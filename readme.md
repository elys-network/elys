# Elys

**Elys** is a blockchain built using Cosmos SDK and CometBFT. It is designed to be a fast, scalable, and secure blockchain that can be used to build decentralized applications.

| Parameter            | Value                                                                    |
| -------------------- | ------------------------------------------------------------------------ |
| Chain Info           | [See network details here](https://github.com/elys-network/networks)     |
| Denomination         | uelys                                                                    |
| Decimals             | 6 (1 elys= 1000000uelys)                                                 |
| Version              | v6.2.0                                                                   |
| MainNet RPC Endpoint | https://rpc.elys.network:443                                             |
| MainNet API Endpoint | https://api.elys.network:443                                             |
| MainNet Explorer     | https://explorer.elys.network ; https://elysscan.io                      |
| TestNet RPC Endpoint | https://rpc.testnet.elys.network:443                                     |
| TestNet API Endpoint | https://api.testnet.elys.network:443                                     |

## Localnet Setup Guide

This guide provides instructions on how to spin up a new localnet using the Elys network for development purposes. Follow these steps to set up your localnet environment.

### Prerequisites

- Make sure you have `git`, `make`, Go environment, and `jq` installed on your machine.

### Getting Started

1. **Clone the Elys Repository**

   First, clone the Elys repository to your local machine:

   ```bash
   git clone https://github.com/elys-network/elys.git
   ```

2. **Build the Binary**

   Navigate into the cloned repository and build the binary using:

   ```bash
   git tag -f v999999.999999.999999 && make install
   ```

   This command will install the `elysd` daemon.

3. **Download the Latest TestNet Snapshot**

   To get the latest TestNet snapshot available for the Elys network, use the following command to download the latest TestNet snapshot that uses the changes from the `main` branch:

   ```bash
   rm -rf ~/.elys && curl -o - -L https://snapshots.elys.network/elys-snapshot-main.tar.lz4 | lz4 -c -d - | tar -x -C ~/
   ```

4. **Spin Up the Localnet**

   Use the command below to start the localnet:

   ```bash
   elysd start
   ```

## Installation

### With Makefile (Recommended)

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
git checkout [version]
```

4. Ensure that you have the necessary dependencies installed. For instance, on Ubuntu you need to install the `make` tool:

```bash
sudo apt-get install --yes make
```

In order to generate proto files, install the dependencies below:

- `buf`
- `clang-format`
- `protoc-gen-go-cosmos-orm`: `go install cosmossdk.io/orm/cmd/protoc-gen-go-cosmos-orm@latest`

Then run the following command:

```bash
make proto
```

5. **Optional**: Use _RocksDB_ instead of _pebbledb_

Ensure that you have RocksDB installed on your machine. On Ubuntu, you can install RocksDB using the following suite of commands:

```bash
# set rocks db version
ROCKSDB_VERSION=8.9.1

# install rocks db dependencies
sudo apt install -y libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev libzstd-dev

# download and extract on /tmp
cd /tmp
wget https://github.com/facebook/rocksdb/archive/refs/tags/v${ROCKSDB_VERSION}.tar.gz
tar -xvf v${ROCKSDB_VERSION}.tar.gz && cd rocksdb-${ROCKSDB_VERSION} || return

# build rocks db
export CXXFLAGS='-Wno-error=deprecated-copy -Wno-error=pessimizing-move -Wno-error=class-memaccess'
make shared_lib

# install rocks db
sudo make install-shared INSTALL_PATH=/usr

# cleanup to save space
rm -rf /tmp/rocksdb-${ROCKSDB_VERSION} /tmp/v${ROCKSDB_VERSION}.tar.gz
```

In order to build the binary with RocksDB, you need to run the following command:

```bash
ROCKSDB=1 make build
```

Note: RocksDB is only required for Linux machines. For macOS, you can continue without installing RocksDB.

When running `ROCKSDB=1 make build`, if you are getting this error:

```bash
elysd: error while loading shared libraries: librocksdb.so.8.9: cannot open shared object file: No such file or directory
```

You might need to set the `LD_LIBRARY_PATH` environment variable to the local library path. You can do this by running the following command:

```bash
export LD_LIBRARY_PATH=/usr/local/lib
```

6. Run the `make build` command to build the binary:

```bash
make build
```

7. The binary will be generated in the `./build` directory. You can run the binary using the following command:

```bash
./build/elysd
```

You can also use the `make install` command to install the binary in the `bin` directory of your `GOPATH`.

</details>

## Validator Guide

The validator guide is accessible [here](./validator.md).

## Architecture

The architecture guide is accessible [here](./architecture.md).

## Release

To release a new version of Elys, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

## Learn more

- [X (formerly Twitter)](https://x.com/elys_network)
- [TestNet Explorer](https://testnet.ping.pub/elys)
- [Developer Chat](https://discord.gg/elysnetwork)
- [Github](https://github.com/elys-network)
