<!--
order: 1
-->

# Concepts

`leveragelp` module provides interface for users to add liquidity in leverage to get more rewards for liquidity providers.

The underlying mechanism is when a user would like to open a leveraged position, module borrows `USDC` from stablestake module and put it as liquidity along with collateral amount received from the user.

To make it secure, the module monitors position value and if it goes down below threshold, it is force closing the position and returning all the borrowed `USDC` back to `stablestake` module along with stacked interest.

# Mechanism

The LeverageLP module enables participants to leverage their liquidity position in the AMM pool by borrowing USDC tokens bonded in the StableStake module. This interaction results in a leveraged position that is tracked and managed within the LeverageLP module.

## Dependencies

The LeverageLP module depends on two other modules:

1. **StableStake**: Manages the bonding of USDC tokens by participants.
2. **AMM**: Facilitates liquidity provision and trading.

## Overview

Participants in the system engage with the modules as follows:

1. **Bonding USDC to StableStake**:

   - Participants bond their USDC tokens to the StableStake module.
   - These bonded tokens become available for borrowing by other participants.

2. **Leveraging Liquidity with LeverageLP**:
   - Participants borrow USDC tokens from the StableStake module using the LeverageLP module.
   - The borrowed tokens are used to enhance their liquidity position in the AMM pool.
   - This results in a debt recorded within the StableStake module.

## Key Equations

To maintain the stability and health of the pool, the following equations are used:

1. **Pool Health**:

   $\text{Pool Health} = \frac{\text{Total Shares} - \text{Leveraged LP Amount}}{\text{Total Shares}}$

   - **Total Shares**: The total number of shares in the pool.
   - **Leveraged LP Amount**: The amount of liquidity added to the pool through leveraged positions.

   The Pool Health metric is monitored to ensure it does not fall below the Safety Factor threshold, which is set at 10%. Lower Pool Health indicates higher risk of leveraged positions being liquidated.

2. **Position Health and Liquidation Mechanism**:

   $\text{Position Health} = \frac{\text{Position Value}}{\text{Borrowed Amount} + \text{Interest Staked} + \text{Interest Paid}}$

   - **Position Value**: The current value of the leveraged position.
   - **Borrowed Amount**: The total amount of USDC borrowed from the StableStake module.
   - **Interest Staked**: The amount of interest accrued on the borrowed amount.
   - **Interest Paid**: The amount of interest paid by the participant.

   If Position Health falls below the Safety Factor threshold (10%), the leveraged LP position is automatically liquidated to protect the underlying liquidity pool.
