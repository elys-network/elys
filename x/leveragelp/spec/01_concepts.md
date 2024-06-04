<!--
order: 1
-->

# Concepts

`leveragelp` module provides interface for users to add liquidity in leverage to get more rewards for liquidity providers.

The underlying mechanism is when a user would like to open a leveraged position, module borrows `USDC` from stablestake module and put it as liquidity along with collateral amount received from the user.

To make it secure, the module monitors position value and if it goes down below threshold, it is force closing the position and returning all the borrowed `USDC` back to `stablestake` module along with stacked interest.
