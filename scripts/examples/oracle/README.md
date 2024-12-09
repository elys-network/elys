# Local network connection

### Start custom bandchain

```
git clone git@github.com:elys-network/bandchain.git
git checkout v2_4_1_customization
go install ./cmd/bandd
sh ./scripts/examples/oracle/start_band.sh
```

### Start elys chain

```bash
elysd init my-node --chain-id elys

# Add a key for alice (if not exists)
elysd keys add alice --keyring-backend test

# Add genesis account
elysd add-genesis-account alice 100000000stake --keyring-backend test

# Generate genesis tx
elysd gentx alice 100000000stake --chain-id elys --keyring-backend test

# Collect genesis txs
elysd collect-gentxs

# Start the chain
elysd start
```

### Start relayer

#### Configure and setup relayer

```bash
hermes config init
hermes keys add --chain band-test --mnemonic-file band_mnemonic.txt
hermes keys add --chain elys --mnemonic-file elys_mnemonic.txt

# Create hermes config file with appropriate settings
# Configure chains, ports, gas settings etc in ~/.hermes/config.toml

# Start the relayer
hermes start
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

### Start local elys chain

```bash
elysd start
```

### Setup relayer between local elys and public band testnet

```bash
# Initialize Hermes config
hermes config init

# Add the chains to your config file (~/.hermes/config.toml):
[[chains]]
id = 'band-laozi-testnet6'
rpc_addr = 'https://rpc.laozi-testnet6.bandchain.org:443'
grpc_addr = 'https://laozi-testnet6.bandchain.org:9090'
websocket_addr = 'wss://rpc.laozi-testnet6.bandchain.org/websocket'
rpc_timeout = '10s'
account_prefix = 'band'
key_name = 'band-relayer'
store_prefix = 'ibc'
gas_price = { price = 0.1, denom = 'uband' }
gas_limit = 5000000

[[chains]]
id = 'elys'
rpc_addr = 'http://localhost:26657'
grpc_addr = 'http://localhost:9090'
websocket_addr = 'ws://localhost:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'elys'
key_name = 'elys-relayer'
store_prefix = 'ibc'
gas_price = { price = 0.0, denom = 'stake' }
gas_limit = 300000

# Add keys for both chains
hermes keys add --chain band-laozi-testnet6 --mnemonic-file band_mnemonic.txt
hermes keys add --chain elys --mnemonic-file elys_mnemonic.txt

# Create the connection
hermes create connection --a-chain elys --b-chain band-laozi-testnet6

# Start the relayer
hermes start
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
