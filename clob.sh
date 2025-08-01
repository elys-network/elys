#!/usr/bin/env bash


TX_RESPONSE_1=$(echo y | elysd tx clob create-market uatom uusdc 0.02 0.01 0.005 0.005 0.005 0.001 0.001 0.5 0.05 0.001 100 --from test1 --chain-id elys-1 --fees 3000uelys  --gas 500000  --keyring-backend="test" -o json)

elysd q clob markets

TX_RESPONSE_1=$(echo y | elysd tx clob deposit 10000000000uusdc --from test1 --chain-id elys-1 --fees 3000uelys  --gas 500000  --keyring-backend="test" -o json)

TX_RESPONSE_2=$(echo y | elysd tx clob deposit 10000000000uusdc --from test2  --chain-id elys-1 --fees 3000uelys  --gas 500000  --keyring-backend="test" -o json)

elysd q clob subaccounts elys1edsl95fuhyx0vlktyz4rsgrp46xsvhle9mjgwv

elysd q clob subaccounts elys1tpcjhfdwenlfuk3ryllyzvs9haavxgx58auqq5

TX_RESPONSE_1=$(echo y | elysd tx clob place-limit-order 1 6.5 1000 limit_buy --isolated-order --from test1 --chain-id elys-1 --fees 1000uelys)

elysd q clob order-book 1 true

TX_RESPONSE_1=$(echo y | elysd tx clob place-limit-order 1 6.4 1000 limit_sell --isolated-order --from test2 --chain-id elys-1 --fees 1000uelys)

elysd q clob order-book 1 false