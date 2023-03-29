# Test band oracle with bootstrapped codebase

https://docs.ignite.com/v0.25.2/kb/band

https://laozi-testnet6.bandchain.org/faucet

rm -rf ~/.ignite/relayer

ignite chain serve --verbose

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

elysd tx oracle coin-rates-data 37 4 3 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id elys -y

elysd query oracle last-coin-rates-id

elysd query oracle coin-rates-result 0

# Local network connection

sh ./x/oracle/start_band.sh

ignite chain serve --verbose

ignite account show abc --address-prefix=band
bandd tx bank send validator band12zyg3xanvupc6upytvsghhlkl8l9cm2rtzn57q 1000000uband --keyring-backend=test --chain-id=band-test -y --broadcast-mode=block --node=http://localhost:26658
bandd query bank balances band12zyg3xanvupc6upytvsghhlkl8l9cm2rtzn57q --node=http://localhost:26658

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

lsof -i :9091

elysd tx oracle coin-rates-data 37 4 3 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id elys -y

elysd query oracle last-coin-rates-id

elysd query oracle coin-rates-result 0
