# Build docker images locally

```
make build-docker build=ebrhf
```

The command to build `elys`, `band`, `go-relayer`, `hermes`, `price-feeder`.

Note: ensure price-feeder repo is cloned under elys repository.

# Build daemon on build directory

The configuration files and genesis files are generated with the binaries built on local and should ensure that correct versions of `bandd`, `elysd`, `hermes`, `price-feeder`, `rly` daemons are put on `build` directory.

### Version table

```
hermes 1.4.0+daad028
rly v2.3.0
price-feeder `4a977ab` (`feature/binance_price_feeder_working` branch)
bandd `79456dc` (`v2_4_1_customization` branch)
```

# Starting dockernet

```sh
make start-docker

# Logs
canceled
canceled
canceled
docker-compose -f ./dockernet/docker-compose.yml down --remove-orphans
[+] Running 6/6
 ⠿ Container dockernet-relayer-band-1  Removed                                                                                                10.2s
 ⠿ Container dockernet-elys3-1         Removed                                                                                                 0.3s
 ⠿ Container dockernet-elys1-1         Removed                                                                                                 0.4s
 ⠿ Container dockernet-elys2-1         Removed                                                                                                 0.3s
 ⠿ Container dockernet-band1-1         Removed                                                                                                 0.5s
 ⠿ Network dockernet_default           Removed                                                                                                 0.1s
Initializing ELYS chain...
Node #1 ID: d031ff99190a741cb0db26259741e792af40d613@elys1:26656
Node #2 ID: 39bfa155ad11ef25e864ed28fb2f3f4e52dffe96@elys2:26656
Node #3 ID: 4091cebb593b7e72ca89feffca8d8bfb9d7d3ee2@elys3:26656
Initializing BAND chain...
Node #1 ID: b28a089a4c34a651f26e745155f9b195db8a6ff5@band1:26656
Starting ELYS chain
[+] Running 4/4
 ⠿ Network dockernet_default    Created
 ⠿ Container dockernet-elys1-1  Started
 ⠿ Container dockernet-elys2-1  Started
 ⠿ Container dockernet-elys3-1  Started
Starting BAND chain
[+] Running 1/1
 ⠿ Container dockernet-band1-1  Started
Waiting for ELYS to start...Done
Waiting for BAND to start...Done
ELYS <> BAND - Adding relayer keys...Done restoring relayer keys
ELYS <> BAND - Creating client, connection, and transfer channel...Done.
[+] Running 1/1
 ⠿ Container dockernet-relayer-band-1  Started
```

# Checking logs from docker

`dockernet/logs` keep the logs from multiple binaries on docker containers.

# Checking the configuration files for dockernet (volumes shared with dockernet)

`dockernet/state` keep the directories for node home directories, price feeder config and relayer configuration.

# Interacting with dockernet

The commands on `dockernet/tests` could be executed to interact with dockernet.

# Testing dockernet environment automatically

```sh
cd dockernet/tests/
sh check.sh

# Logs
+ source ./assert.sh
++ command -v tput
++ tty -s
+++ tput setaf 1
++ RED=''
+++ tput setaf 2
++ GREEN=''
+++ tput setaf 5
++ MAGENTA=''
+++ tput sgr0
++ NORMAL=''
+++ tput bold
++ BOLD=''
++ elysd query oracle show-price BTC --source=band
+ band_price='price:
  asset: BTC
  price: "0.000000000000490441"
  provider: automation
  source: band
  timestamp: "1683195717"'
+ assert_contain 'price:
  asset: BTC
  price: "0.000000000000490441"
  provider: automation
  source: band
  timestamp: "1683195717"' 'asset: BTC' 'band price check failure'
+ local 'haystack=price:
  asset: BTC
  price: "0.000000000000490441"
  provider: automation
  source: band
  timestamp: "1683195717"'
+ local 'needle=asset: BTC'
+ local 'msg=band price check failure'
+ '[' -z x ']'
+ '[' -z '' ']'
+ return 0
++ elysd query oracle show-price BTC --source=binance
+ elys_price='price:
  asset: BTC
  price: "29159.631493739183154355"
  provider: elys1mxk8wmns33vs6yynsaeud2k97xkl5dqlkjv3j9
  source: binance
  timestamp: "1683195535"'
+ assert_contain 'price:
  asset: BTC
  price: "29159.631493739183154355"
  provider: elys1mxk8wmns33vs6yynsaeud2k97xkl5dqlkjv3j9
  source: binance
  timestamp: "1683195535"' 'asset: BTC' 'elys price check failure'
+ local 'haystack=price:
  asset: BTC
  price: "29159.631493739183154355"
  provider: elys1mxk8wmns33vs6yynsaeud2k97xkl5dqlkjv3j9
  source: binance
  timestamp: "1683195535"'
+ local 'needle=asset: BTC'
+ local 'msg=elys price check failure'
+ '[' -z x ']'
+ '[' -z '' ']'
+ return 0
+ echo 'All passed'
All passed
```
