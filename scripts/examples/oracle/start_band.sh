#!/bin/sh

rm -rf $HOME/.band/

cd $HOME

chain_id=band-test

echo "figure web rescue rice quantum sustain alert citizen woman cable wasp eyebrow monster teach hockey giant monitor hero oblige picnic ball never lamp distance" > mnt.txt;

bandd init --chain-id=$chain_id band-test --home=$HOME/.band
bandd keys add validator --keyring-backend=test --home=$HOME/.band
bandd keys add user1 --recover --keyring-backend=test < mnt.txt;
bandd add-genesis-account $(bandd keys show validator -a --keyring-backend=test --home=$HOME/.band) 100000000000uband --home=$HOME/.band
bandd add-genesis-account $(bandd keys show user1 -a --keyring-backend=test --home=$HOME/.band) 10000000uband --home=$HOME/.band
bandd add-genesis-account band12zyg3xanvupc6upytvsghhlkl8l9cm2rtzn57q 10000000uband --home=$HOME/.band
bandd gentx validator 500000000uband --keyring-backend=test --home=$HOME/.band --chain-id=$chain_id
bandd collect-gentxs --home=$HOME/.band

sed -i '' 's#"localhost:6060"#"localhost:6061"#g' $HOME/.band/config/config.toml
sed -i '' 's#"0.0.0.0:9091"#"0.0.0.0:9095"#g' $HOME/.band/config/app.toml

bandd start --home=$HOME/.band --rpc.laddr tcp://localhost:26658 \
	--grpc.address localhost:9094 \
    --address tcp://localhost:26655 \
    --p2p.laddr tcp://localhost:26656


