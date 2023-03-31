#!/usr/bin/env bash

# query params
elysd query oracle coin-rates-result
elysd query oracle last-coin-rates-id 1
elysd query oracle list-asset-info
elysd query oracle show-asset-info satoshi
elysd query oracle params
elysd query gov proposals

elysd tx oracle coin-rates-data
elysd tx oracle add-asset-info-proposal satoshi BTC BTC "" "" --title="title" --description="description" --deposit="10000000000000000000aelys" --from=alice --chain-id=elys --broadcast-mode=block --yes
elysd tx oracle remove-asset-info-proposal satoshi --title="title" --description="description" --deposit="10000000000000000000aelys" --from=alice --chain-id=elys --broadcast-mode=block --yes
elysd tx oracle create-asset-info satoshi BTC BTC "" "" --from=alice --chain-id=elys --broadcast-mode=block --yes
elysd tx oracle update-asset-info satoshi BTC BTC "" "" --from=alice --chain-id=elys --broadcast-mode=block --yes
elysd tx oracle delete-asset-info satoshi --from=alice --chain-id=elys --broadcast-mode=block --yes

elysd tx gov vote 1 yes --from=alice --chain-id=elys --broadcast-mode=block --yes
elysd tx gov vote 2 yes --from=alice --chain-id=elys --broadcast-mode=block --yes