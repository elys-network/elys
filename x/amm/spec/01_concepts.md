<!--
order: 1
-->

# Concepts

The goal of `oracle based pool` is to offer swap and perpetual traders a low slippage UX (arbitrage free price) despite the amount of liquidity. It combines `balancer` slippage and `oracle` based weight model to provide a good price while protecting liquidity from arbitrageurs.

New model aims to increase fund utilization ratio and to provide more APR to liquidity providers.

To keep the asset weight at target weight, weight breaking fee is introduced to be spent when weight is broken more than `threshold`.
