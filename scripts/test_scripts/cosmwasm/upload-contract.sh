export ELYSD_NODE="tcp://localhost:26657"

FLAGS="--from validator --gas 10000000 --keyring-backend=test --chain-id=elysicstestnet-1 --output=json --yes"

elysd tx wasm submit-proposal wasm-store ./scripts/test_scripts/cosmwasm/cw_template.wasm --title 'first auth contract' --summary 'hehe'  --deposit 100000000uelys $FLAGS
sleep 2

elysd tx gov vote 32 Yes $FLAGS
sleep 60

txhash=$(elysd tx wasm submit-proposal instantiate-contract 2 '{"count":0}' --label 'hello-test' --title 'test' --summary 'test-sum' --no-admin --deposit 100000000uelys $FLAGS | jq -r .txhash) && echo $txhash
echo "txhash for instantiate contract: $txhash"
sleep 2
elysd tx gov vote 34 Yes $FLAGS
echo "wait for proposal to pass"