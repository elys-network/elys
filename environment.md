# Introduction

In this document, we're going to see how to create a new environment from scratch. We will demonstrate it setting up 2 different nodes with the Elys binary on its version **0.30.0**.

## Nodes

1.  Mallorca: Represents 60% of the staked tokens.
2.  Fuji: Represents 40% of the staked tokens.

## Data for this example:

- chain-id: `elysdevnet-1`
- Mallorca moniker: `mallorca`
- Fuji moniker: `fuji`

## Install Elys on (Mallorca and Fuji)

### Install the node

> You have to make this step with Mallorca and Fuji Servers

```bash
cd $HOME || return
rm -rf $HOME/elys
git clone https://github.com/elys-network/elys.git
cd $HOME/elys || return
git checkout v0.30.0

make install

elysd config keyring-backend os
elysd config chain-id elysdevnet-1

#Replace the moniker en each server with mallorca o fuji:
#Examples:
#elysd init "fuji" --chain-id elysdevnet-1
#elysd init "mallorca" --chain-id elysdevnet-1
elysd init "Your Moniker" --chain-id elysdevnet-1

APP_TOML="~/.elys/config/app.toml"
sed -i "s/^app-db-backend *=.*/app-db-backend = \"pebbledb\"/" $APP_TOML
sed -i 's|^pruning *=.*|pruning = "custom"|g' $APP_TOML
sed -i 's|^pruning-keep-recent  *=.*|pruning-keep-recent = "100"|g' $APP_TOML
sed -i 's|^pruning-keep-every *=.*|pruning-keep-every = "0"|g' $APP_TOML
sed -i 's|^pruning-interval *=.*|pruning-interval = 19|g' $APP_TOML
CONFIG_TOML="~/.elys/config/config.toml"
external_address=$(wget -qO- eth0.me)
sed -i.bak -e "s/^external_address *=.*/external_address = \"$external_address:26656\"/" $CONFIG_TOML
sed -i 's|^minimum-gas-prices *=.*|minimum-gas-prices = "0.0018ibc/2180E84E20F5679FCC760D8C165B60F42065DEF7F46A72B447CFF1B7DC6C0A65,0.00025ibc/E2D2F6ADCC68AA3384B2F5DFACCA437923D137C14E86FB8A10207CF3BED0C8D4,0.00025uelys"|g' $CONFIG_TOML
sed -i 's|^prometheus *=.*|prometheus = true|' $CONFIG_TOML
sed -i -e "s/^filter_peers *=.*/filter_peers = \"true\"/" $CONFIG_TOML
```

### RPC server extra configuration

Edit your config.toml and change `laddr` and `cors_allowed_origins` like these:

```
[rpc]
laddr = "tcp://0.0.0.0:26657"
cors_allowed_origins = ["*"]
```

## Genesis creation

### Initial Balances and Initial Stakes

**Mallorca Server**

```bash
elysd keys add mallorca
elysd add-genesis-account {mallorca-address} 1000000000uelys
elysd gentx mallorca 60000000uelys \
    --keyring-backend file --keyring-dir /home/ubuntu/.elys/keys \
    --account-number 0 --sequence 0 \
    --pubkey=$(elysd tendermint show-validator) \
    --chain-id elysdevnet-1 \
    --gas 1000000 \
    --gas-prices 0.1uelys
```

Share the genesis file to Fuji server and continue:

**Fuji Server**

```bash
elysd keys add fuji
elysd add-genesis-account {fuji-address} 1000000000uelys
elysd gentx fuji 40000000uelys \
    --keyring-backend file --keyring-dir /home/ubuntu/.elys/keys \
    --account-number 0 --sequence 0 \
    --pubkey=$(elysd tendermint show-validator) \
    --chain-id elysdevnet-1 \
    --gas 1000000 \
    --gas-prices 0.1uelys
```

## Genesis assembly

With the two initial staking transactions created. Mallorca have to include both of them in the genesis:

1. Move the fuji tx in `.elys/config/gentx/gentx-*` to mallorca: `.elys/config/gentx/`
2. Call to `elysd collect-gentxs` from Mallorca.
3. As an added precaution, confirm that it is a valid genesis: `elysd validate-genesis`
   It should return:
   ` File at /root/.checkers/config/genesis.json is a valid genesis file`
4. Edit your genesis.json and change every reference from `"stake"` denom to `"uelys"`.

## Genesis distribution

All the nodes that will start the chain need the final version of the genesis so please share again with Fuji.

## Network preparation

### Setup Fuji

1. `elysd tendermint show-node-id`
2. This returns something like: f2673103417334a839f5c20096909c3023ba4903
3. Your node identification is: f2673103417334a839f5c20096909c3023ba4903@your-public-ip:26656

Eventually, in fuji server, you should have **config/config.toml**:

```bash
seeds = "7009cc51174dce87c31f537fe8fed906349a27f4@ip-mallorca:26656"
persistent_peers = "f2673103417334a839f5c20096909c3023ba4903@fuji-ip:26656"
private_peer_ids = "f2673103417334a839f5c20096909c3023ba4903"
```

### Setup Mallorca

1. `elysd tendermint show-node-id`
2. This returns something like: 009cc51174dce87c31f537fe8fed906349a27f4
3. Your node identification is: 009cc51174dce87c31f537fe8fed906349a27f4@your-public-ip:26656

Eventually, in mallorca server, you should have **config/config.toml**:

```bash
seeds = "f2673103417334a839f5c20096909c3023ba49034@fuji:26656"
persistent_peers = "7009cc51174dce87c31f537fe8fed906349a27f43@mallorca-ip:26656"
private_peer_ids = "7009cc51174dce87c31f537fe8fed906349a27f4"
```

## Cosmovisor

We will be setting up Cosmovisor.
Cosmovisor is a process manager for Cosmos SDK application binaries that automates application binary switch at chain upgrades.

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.4.0
mkdir -p ~/.elys/cosmovisor/genesis/bin
mkdir -p ~/.elys/cosmovisor/upgrades
cp ~/go/bin/elysd ~/.elys/cosmovisor/genesis/bin
```

## Starting the chain as a Service

```bash
sudo tee /etc/systemd/system/elysd.service > /dev/null << EOF
[Unit]
Description=Elys Node
After=network-online.target
[Service]
User=$USER
ExecStart=$(which cosmovisor) run start
Restart=on-failure
RestartSec=3
LimitNOFILE=10000
Environment="DAEMON_NAME=elysd"
Environment="DAEMON_HOME=$HOME/.elys"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="UNSAFE_SKIP_BACKUP=true"
[Install]
WantedBy=multi-user.target
EOF
```

[Reference](https://services.stake-town.com/home/testnet/elys/installation)

## Enable and start service

```bash
sudo systemctl daemon-reload
sudo systemctl enable elysd
sudo systemctl start elysd
```

## Stopping the service

`sudo systemctl stop elysd`

## Watching logs

`sudo journalctl -u elysd -f -o cat`
