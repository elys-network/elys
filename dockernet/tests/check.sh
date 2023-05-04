#!/bin/bash
set -eux

source "./assert.sh"

# elysd query oracle list-price-feeder
# elysd query oracle list-price
# elysd query epochs epoch-infos
# elysd query oracle params
# elysd query oracle list-asset-info

# elysd query oracle show-price BTC --source=band
# elysd query oracle show-price BTC --source=elys

band_price=$(elysd query oracle show-price BTC --source=band)
assert_contain "$band_price" "asset: BTC" "band price check failure"

elys_price=$(elysd query oracle show-price BTC --source=binance)
assert_contain "$elys_price" "asset: BTC" "elys price check failure"
echo "All passed"
