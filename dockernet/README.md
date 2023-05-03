# Build docker images locally

```
make build-docker build=ebrhf
```

The command to build `elys`, `band`, `go-relayer`, `hermes`, `price-feeder`.

Note: ensure price-feeder repo is cloned under elys repository.

# Build daemon on build directory

The configuration files and genesis files are generated with the binaries built on local and should ensure that correct versions of `bandd`, `elysd`, `hermes`, `price-feeder`, `rly` daemons are put on `build` directory.

# Starting dockernet

```
make start-docker
```

# Checking logs from docker

`dockernet/logs` keep the logs from multiple binaries on docker containers.

# Checking the configuration files for dockernet (volumes shared with dockernet)

`dockernet/state` keep the directories for node home directories, price feeder config and relayer configuration.

# Interacting with dockernet

The commands on `dockernet/tests` could be executed to interact with dockernet.
