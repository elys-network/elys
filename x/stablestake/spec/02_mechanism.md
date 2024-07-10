<!--
order: 2
-->

# Mechanism

## Overview

The StableStake module facilitates leveraged liquidity provision (LP) positions by allowing users to borrow funds against their collateral. This process involves interactions between the user, a stable stake pool, and an Automated Market Maker (AMM) pool. The module ensures efficient and secure fund management, enabling users to maximize their liquidity provisioning potential. Additionally, users can bond and unbond their USDC to earn rewards from the interest paid by ongoing leveraged LP positions.

## Components

### Symbols and Definitions

- **Collateral**: The amount of assets provided by the user as security for the borrowed funds.
- **Borrowed Amount**: The funds provided by the stable stake pool for leverage.
- **Position Address**: The designated address handling position-specific funds.
- **AMM Pool**: The decentralized exchange pool where liquidity is provided.
- **Stable Stake Pool**: The source of borrowed funds for leverage.

### Key Entities

- **User**: The entity initiating the opening and closing of leveraged LP positions.
- **Position Owner**: The user who owns the collateral and manages the position.
- **StableStake Module**: The system managing the stable stake pool and the borrowed funds.
- **AMM Module**: The system managing the liquidity provision and trading in the AMM pool.
- **LeverageLP Module**: The module facilitating leveraged liquidity provision positions.

## Workflow

### 1. Opening a Leveraged LP Position

#### Step 1: Initiate Position Opening

- **Action**: User decides to open a leveraged LP position.
- **Trigger**: User sends a request to open a position.

#### Step 2: Send Collateral Amount

- **Action**: The user sends a collateral amount from their address to the position address.
- **Details**:
  - **Source**: User's wallet (position owner)
  - **Destination**: Position address
  - **Amount**: Collateral specified by the user

#### Step 3: Send Borrowed Amount

- **Action**: The system sends the borrowed amount from the stable stake pool to the position address.
- **Details**:
  - **Source**: Stable stake pool
  - **Destination**: Position address
  - **Amount**: Borrowed funds required for leverage

#### Step 4: Join AMM Pool

- **Action**: The system joins the AMM pool with the combined collateral and borrowed amounts.
- **Details**:
  - **Source**: Position address
  - **Destination**: AMM pool
  - **Amount**: Collateral + Borrowed funds

### 2. Closing a Leveraged LP Position

#### Step 1: Initiate Position Closing

- **Action**: User decides to close the leveraged LP position.
- **Trigger**: User sends a request to close the position.

#### Step 2: Exit AMM Pool

- **Action**: The system exits the AMM pool to retrieve both the collateral and borrowed amounts towards the position address.
- **Details**:
  - **Source**: AMM pool
  - **Destination**: Position address
  - **Amount**: Collateral + Borrowed funds

#### Step 3: Send Borrowed and Debts Amounts

- **Action**: The system sends both borrowed and debt amounts from the position address to the stable stake pool.
- **Details**:
  - **Source**: Position address
  - **Destination**: Stable stake pool
  - **Amount**: Borrowed funds + Accrued interest/debts

#### Step 4: Return Remaining Amount

- **Action**: The system sends the remaining amount back to the position owner.
- **Details**:
  - **Source**: Position address
  - **Destination**: User's wallet (position owner)
  - **Amount**: Remaining collateral after settling the debt

## Sequence of Actions

### Bonding and Unbonding USDC

- Bonding: Users can bond their USDC to the stable stake pool to earn rewards. The rewards are generated from the interest paid by users on their leveraged LP positions.
- Unbonding: Users can unbond their USDC from the stable stake pool, allowing them to retrieve their tokens along with any accrued rewards.

### Opening Position

1. User sends a request to open a leveraged LP position.
2. Collateral is transferred from the user's wallet to the position address.
3. Borrowed funds are transferred from the stable stake pool to the position address.
4. The position address joins the AMM pool with the total funds.

### Closing Position

1. User sends a request to close the leveraged LP position.
2. The position address exits the AMM pool, retrieving the collateral and borrowed amounts.
3. Borrowed and debt amounts are transferred from the position address to the stable stake pool.
4. Remaining collateral is sent back to the user's wallet.
