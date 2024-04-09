# How to use a testnet snapshot in localnet?

## High Stakes testnet snapshots

You can find the latest High Stakes testnet snapshot [here](https://tools.highstakes.ch/files/elys.tar.gz).

```
make install
go run ./scripts/upgrade-assure/... https://tools.highstakes.ch/files/elys.tar.gz ~/go/bin/elysd ~/go/bin/elysd --skip-proposal
```

## Stake Town testnet snapshots

You can find the latest Stake Town testnet snapshot [here](https://snapshots-testnet.stake-town.com/elys/elystestnet-1_latest.tar.lz4).

```
make install
go run ./scripts/upgrade-assure/... https://snapshots-testnet.stake-town.com/elys/elystestnet-1_latest.tar.lz4 ~/go/bin/elysd ~/go/bin/elysd --skip-proposal
```

## Polkachu testnet snapshots

You can find the latest Polkachu testnet snapshot [here](https://polkachu.com/testnets/elys/snapshots).

```
make install
go run ./scripts/upgrade-assure/... https://snapshots.polkachu.com/testnet-snapshots/elys/elys_5724942.tar.lz4 ~/go/bin/elysd ~/go/bin/elysd --skip-proposal
```

## AviaOne testnet snapshots

You can find the latest AviaOne testnet snapshot [here](https://aviaone.com/blockchains-service/elystestnet-1-elys.html#8).

```
make install
go run ./scripts/upgrade-assure/... https://services.elystestnet-1.elys.aviaone.com/elystestnet-1_2024-03-06.tar.gz ~/go/bin/elysd ~/go/bin/elysd --skip-proposal
```

# How can I perform a test with a version upgrade that involves extensive changes to data structures?

```
git checkout v0.29.30
make install
cp -a ~/go/bin/elysd /tmp/elysd-v0.29.30
```

```
go run ./scripts/upgrade-assure/... --home /tmp/elys --home2 /tmp/elysd2 https://tools.highstakes.ch/files/elys.tar.gz /tmp/elysd-v0.29.30 /tmp/elysd-v0.29.31 --skip-node-start
```

```
git checkout v0.29.31
make install
cp -a ~/go/bin/elysd /tmp/elysd-v0.29.31
```

```
go run ./scripts/upgrade-assure/... --home /tmp/elys --home2 /tmp/elysd2 https://tools.highstakes.ch/files/elys.tar.gz /tmp/elysd-v0.29.30 /tmp/elysd-v0.29.31 --skip-snapshot --skip-chain-init
```
