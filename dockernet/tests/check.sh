#!/bin/bash

elysd query oracle list-price-feeder
elysd query oracle list-price
elysd query epochs epoch-infos
elysd query oracle params
elysd query oracle list-asset-info

elysd query oracle show-price BTC --source=band
elysd query oracle show-price BTC --source=elys
