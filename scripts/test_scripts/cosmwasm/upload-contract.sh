export ELYSD_NODE="tcp://localhost:26657"

FLAGS="--from validator --gas 10000000 --keyring-backend=test --chain-id=elysicstestnet-1 --output=json --yes"

elysd tx wasm submit-proposal wasm-store ./test_scripts/cosmwasm/cw_template.wasm --title 'first auth contract' --summary 'hehe'  --deposit 100000000uelys $FLAGS
sleep 2

elysd tx gov vote 1 Yes $FLAGS
sleep 60

txhash=$(elysd tx wasm submit-proposal instantiate-contract 1 '{"count":0}' --label 'hello-test' --title 'test' --summary 'test-sum' --authority elys1g3qnq7apxv964cqj0hza0pnwsw3q920lcc5lyg --no-admin $FLAGS | jq -r .txhash) && echo $txhash
elysd tx gov vote 2 Yes $FLAGS
sleep 60

# addr=$(elysd q tx $txhash --output=json | jq -r .logs[0].events[2].attributes[0].value) && echo $addr
# sleep 2

# elysd q wasm contract $addr