
# Specification: Incentivized Liquidation System

## Overview
The Incentivized Liquidation System ensures the stability and health of the leverage market by allowing users and bots to monitor and liquidate undercollateralized positions and postions that met a certain price(stop loss). Liquidators are rewarded for their efforts, creating an incentivized mechanism to maintain market health.

## Scalability Advantage
Unlike a system that requires traversing and evaluating all positions in each block—which is not scalable—the incentivized liquidation approach delegates the monitoring and liquidation tasks to users and bots. This significantly reduces the computational load on the blockchain, making the system scalable and capable of handling a large number of positions efficiently.

## Incentivized Liquidation by Users and Bots

Users and bots can monitor positions and trigger liquidations when necessary. They are rewarded for their actions, ensuring continuous monitoring and prompt liquidation of undercollateralized positions and stop loss.

### Pseudo Code

```plaintext
function IncentivizedCheckAndLiquidate(ctx, positionId, liquidator):
    position = getPosition(ctx, positionId)
    health = calculateHealth(position)
    
    if health < LIQUIDATION_THRESHOLD || lp_price <= stop_loss_price::
        executeLiquidation(ctx, position)
        rewardLiquidator(ctx, liquidator)
```

## Workflow

1. **Monitoring**: Users and bots monitor the health of positions.
2. **Triggering Liquidation**: If a position's health falls below the threshold, users or bots call the `IncentivizedCheckAndLiquidate` function.
3. **Executing Liquidation**: The system executes the liquidation, sells the assets, covers the debt, and marks the position as liquidated.
4. **Rewarding Liquidators**: The liquidator is rewarded for their action.

## Note
If a user reports an incorrect position for liquidation, they will incur the extra fees associated with the transaction.

## Incentives
1. ***Dynamic Reward Mechanism***: The liquidator's reward is dynamically calculated based on the safety factor of the liquidated position. This creates a direct correlation between the risk taken by the liquidator and the reward received.
2. ***Safety Factor Threshold***: For positions liquidated with a safety factor slightly above 1 (e.g., 1.05), the liquidator is entitled to the excess collateral. Specifically, if a position's safety factor is 1.05 at the time of liquidation, the liquidator receives the additional 5% of the position's value as a reward.
3. ***Internal Bots***: In addition to allowing external users and bots to participate in liquidations, we will also run our own bots. These internal bots will serve as a safeguard to ensure that no undercollateralized position goes unnoticed or unliquidated, providing an additional layer of security for the system.

This incentive structure ensures that liquidators are adequately compensated for their role in maintaining market stability. The dynamic and competitive nature of the rewards encourages continuous monitoring and swift action, contributing to the overall health of the leverage market.
