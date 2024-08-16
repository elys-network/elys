# Specification: Incentivized Liquidation System v2

## Overview

The Incentivized Liquidation System is designed to maintain the stability and health of the leverage market by allowing users and bots to monitor and liquidate undercollateralized positions or those that have met a specific price condition (stop loss). Liquidators are rewarded for their contributions, creating an incentivized mechanism to ensure market stability and timely intervention.

## Scalability Advantage

Unlike systems that require evaluating all positions at each block—a method that is computationally intensive and not scalable—the incentivized liquidation system delegates the tasks of monitoring and liquidation to users and bots. This delegation significantly reduces the computational burden on the blockchain, making the system scalable and capable of efficiently managing a large number of positions.

## Incentivized Liquidation by Users and Bots

Users and bots can actively monitor positions and trigger liquidations when certain criteria are met. This decentralized approach ensures continuous oversight and prompt liquidation of undercollateralized positions or those nearing the stop loss threshold.

### Pseudo Code

```plaintext
function IncentivizedCheckAndLiquidate(ctx, positionId, liquidator):
    position = getPosition(ctx, positionId)
    health = calculateHealth(position)

    if health < LIQUIDATION_THRESHOLD or lp_price <= stop_loss_price:
        executeLiquidation(ctx, position)
        rewardLiquidator(ctx, liquidator)
```

## Workflow

1. **Monitoring**: Users and bots continuously monitor the health of positions.
2. **Triggering Liquidation**: When a position's health drops below the liquidation threshold, or the stop loss price is reached, users or bots initiate the `IncentivizedCheckAndLiquidate` function.
3. **Executing Liquidation**: The system processes the liquidation, selling the assets, covering the debt, and marking the position as liquidated.
4. **Rewarding Liquidators**: The liquidator is rewarded for successfully closing the position.

## Penalty Program

- **Submission Requirement**: Users or bots must submit a `closePositions` message with a list of positions to liquidate or close, along with a minimum required deposit amount, as specified by a governance parameter.
- **Penalty Condition**: If any position in the submitted list is healthy or has not met the stop loss criteria, the deposit is forfeited and retained by the system.
- **Deposit Return**: If all positions in the list meet the criteria for liquidation or closure, the deposit is returned to the submitting address.

## Margin of Error

To account for brief fluctuations in position health or stop loss criteria, a margin of error is introduced:

- **Margin Error Rate**: Set at 1%.
- **Application**:
  - **Position Health**: If the position health is within 1% above the liquidation threshold, the position will not be liquidated, and the deposit will be returned without penalty.
  - **Stop Loss**: If the stop loss price is within 1% above the current market price, the position will not be closed, and the deposit is refunded without penalty.

## Incentive Program

Liquidators are incentivized through rewards when they successfully liquidate or close positions:

1. **Liquidation**:

   - **Safety Factor**: A safety factor of 5% is used to determine when liquidation is necessary.
   - **Reward**: The liquidator receives the remaining value of the liquidated position after deducting the community fund rate.

2. **Stop Loss Closure**:
   - **Condition**: The position is closed when the stop loss price is approached (at stop loss price + 1%).
   - **Reward**: The liquidator receives the 1% difference in the position value as an incentive.

## Arbitrage Opportunities

The system allows liquidators to engage in arbitrage by exploiting differences between actual position values and predefined thresholds, further enhancing the incentives for participation.

## Note

Submitting an incorrect position for liquidation incurs penalties, including the loss of the deposit amount as well as the extra fees associated with the transaction. This ensures that only valid and necessary liquidations are executed, maintaining the integrity and efficiency of the system.
