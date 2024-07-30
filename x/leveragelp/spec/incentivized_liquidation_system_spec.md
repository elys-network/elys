
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

## Incentive value (TODO)
- Fixed vs proportional
    - Fixed Incentive: A fixed incentive means that liquidators receive a predetermined reward amount for each successful liquidation, regardless of the size of the position.
    - Proportional Incentive with Minimum Threshold: A proportional incentive means that the reward is a percentage of the liquidated position's assets. To ensure fairness and adequate motivation, a minimum reward threshold is set.
