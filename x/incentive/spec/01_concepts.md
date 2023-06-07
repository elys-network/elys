<!--
order: 1
-->

# Concepts

The Incentive module is designed to reward the ecosystem participants including liquidity providers (LPs), Elys stakers, and Eden committers. We have 2 kinds of rewards in Elys ecosystem - inflationary reward and non-inflationary rewards. 

1. Inflationary rewards
- Given the amount in each epoch, incentive module distribute it to stakers and LPs by increasing uncommitted token balances periodically. It utilizes the commitment module and updates accounting for both Eden tokens and Eden-Boost tokens of each Elys staker and LP.

2. Non-inflationary rewards
- Fees that are collected from amm, margin and transaction gas fees.

3. Functions
- Distribute inflationary rewards and non-inflationary rewards to stakers and LPs.
- It distributes Eden and Eden boost tokens as inflationary rewards.
- It distributes only USDC as non-inflationary rewards. Fees collected from different parts will be converted into USDC using amm module.
- It funds community pool based on community pool tax.
