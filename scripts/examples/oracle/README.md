# Local network connection

### Start custom bandchain

```
git clone git@github.com:elys-network/bandchain.git
git checkout v2_4_1_customization
go install ./cmd/bandd
sh ./scripts/examples/oracle/start_band.sh
```

### Start elys chain

```
ignite chain init
ignite chain serve --verbose
```

### Start relayer

#### Deposit balance to bandchain relayer account

```
ignite account show abc --address-prefix=band
bandd tx bank send validator band12zyg3xanvupc6upytvsghhlkl8l9cm2rtzn57q 1000000uband --keyring-backend=test --chain-id=band-test -y --broadcast-mode=block --node=http://localhost:26658
bandd query bank balances band12zyg3xanvupc6upytvsghhlkl8l9cm2rtzn57q --node=http://localhost:26658
```

#### Configure and setup relayer

```
rm -rf ~/.ignite/relayer

ignite relayer configure -a \
--source-rpc "http://localhost:26658" \
--source-port "oracle" \
--source-gasprice "0.1uband" \
--source-gaslimit 5000000 \
--source-prefix "band" \
--source-version "bandchain-1" \
--target-rpc "http://localhost:26657" \
--target-faucet "http://localhost:4500" \
--target-port "oracle" \
--target-gasprice "0.0stake" \
--target-gaslimit 300000 \
--target-prefix "elys" \
--target-version "bandchain-1"

ignite relayer connect
```

### Send price request to bandchain

```
elysd tx oracle request-band-price 37 1 1 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id elys -y
```

### Check request id created on bandchain received on Elys

```
elysd query oracle last-band-request-id

request_id: "2"
```

### Check response from bandchain for price request

```
elysd query oracle band-price-result 1

result:
  rates:
  - "489540"
  - "979081"
  - "1468621"
  - "1958162"
  - "2447702"
```

# Test band oracle with public testnet

https://docs.ignite.com/v0.25.2/kb/band

### Start local elys chain

ignite chain serve --verbose

### Setup relayer between local elys and public band testnet

```
rm -rf ~/.ignite/relayer
ignite relayer configure -a \
--source-rpc "https://rpc.laozi-testnet6.bandchain.org:443" \
--source-faucet "https://laozi-testnet6.bandchain.org/faucet" \
--source-port "oracle" \
--source-gasprice "0.1uband" \
--source-gaslimit 5000000 \
--source-prefix "band" \
--source-version "bandchain-1" \
--target-rpc "http://localhost:26657" \
--target-faucet "http://localhost:4500" \
--target-port "oracle" \
--target-gasprice "0.0stake" \
--target-gaslimit 300000 \
--target-prefix "elys" \
--target-version "bandchain-1"

ignite relayer connect
```

### Send price request to bandchain

```
elysd tx oracle request-band-price 37 4 3 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id elys -y
```

### Check request id created on bandchain received on Elys

```
elysd query oracle last-band-request-id
```

### Check response from bandchain for price request

```
elysd query oracle band-price-result 0
```
