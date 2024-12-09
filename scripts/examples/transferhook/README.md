# Local network connection

### Start mars chain

ignite chain serve --verbose --reset-once -c ./scripts/examples/transferhook/mars.yml

### Start elys chain

ignite chain serve --verbose --reset-once

### Start relayer

#### Configure and setup hermes

```sh
hermes --config ./scripts/examples/transferhook/config.toml keys delete --chain marstestnet-1 --all
echo "attend copy blast obey agent clinic monkey blur doctor sibling impact stomach judge rubber actress forest wage silent sick key divide exotic junk velvet" > mars.txt
hermes --config ./scripts/examples/transferhook/config.toml keys add --chain marstestnet-1 --mnemonic-file=mars.txt

hermes --config ./scripts/examples/transferhook/config.toml keys delete --chain elysicstestnet-1 --all
echo "olympic slide park figure frost benefit deer reform fly pull price airport submit monitor silk insect uphold convince pupil project ignore roof warfare slight" > elys.txt
hermes --config ./scripts/examples/transferhook/config.toml keys add --chain elysicstestnet-1 --mnemonic-file=elys.txt

hermes --config ./scripts/examples/transferhook/config.toml create connection --a-chain elysicstestnet-1 --b-chain marstestnet-1

hermes --config ./scripts/examples/transferhook/config.toml create channel --a-chain elysicstestnet-1 --a-port transfer --b-port transfer --a-connection connection-0

hermes --config ./scripts/examples/transferhook/config.toml start
```

#### Try to send with memo

```sh
elysd keys add acc1 --keyring-backend=test
ACC1=elys1jfmgwygyf9u3sx6xm36xlxf85v977k5c7l8w3q

elysd tx ibc-transfer transfer transfer channel-0 $ACC1 10000uatom --chain-id=elysicstestnet-1 --from=treasury --keyring-backend=test -y --node=http://localhost:26657
elysd query bank balances $ACC1 --node=http://localhost:26659

ATOM_IBC=ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2
elysd tx ibc-transfer transfer transfer channel-0 $ACC1 100$ATOM_IBC --chain-id=marstestnet-1 --from=acc1 --keyring-backend=test -y --node=http://localhost:26659

elysd tx amm create-pool 10uatom,10uusdt 10000uatom,10000uusdt --swap-fee=0.00 --exit-fee=0.00 --use-oracle=false --from=treasury --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000

# swap with transferhook
elysd tx ibc-transfer transfer transfer channel-0 '{"transferhook":{"receiver":"'$ACC1'","amm":{"action":"Swap","routes":[{"pool_id":1,"token_out_denom": "uusdt"}]}}}' 100$ATOM_IBC --chain-id=marstestnet-1 --from=acc1 --keyring-backend=test -y --node=http://localhost:26659
elysd query bank balances $ACC1 --node=http://localhost:26657

elysd tx amm swap-exact-amount-in 10uatom 1 1 uusdt --from=acc1 --keyring-backend=test --chain-id=elysicstestnet-1 --yes --gas=1000000
elysd query bank balances $ACC1 --node=http://localhost:26657
```
