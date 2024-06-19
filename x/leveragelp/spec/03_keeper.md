<!--
order: 3
-->

# Keeper

## Position Management

The `leveragelp` module's keeper handles the opening, updating, and closing of leveraged positions. It ensures that positions are correctly managed, leverage limits are respected, and health checks are performed to avoid liquidation risks.

### Opening a Position

The `Open` function allows a user to open a leveraged position by providing collateral and specifying leverage and a stop-loss price.

### Closing a Position

The `Close` function allows a user to close a leveraged position, either partially or fully, by specifying the amount of LP tokens to withdraw.

### Liquidation and Health Checks

The `BeginBlocker` function performs regular health checks on all positions and liquidates unhealthy positions if necessary.

### Positions indexing

Two indexings are available to

- sort positions by risk level per pool
- sort positions by stopLossPrice per pool

## BeginBlocker

The `BeginBlocker` function is called at the beginning of each block to perform necessary updates and maintenance for the `leveragelp` module. It processes health checks and liquidates unhealthy positions if needed.

- To avoid time complexity in liquidation flow, it's iterating till it meets the first healthy position.
- To avoid time complexity in stop loss close position, it's iterating till it meets the first acceptable stop loss unclose position.

### LiquidatePositionIfUnhealthy

The `LiquidatePositionIfUnhealthy` function checks the health of a position and liquidates it if it is unhealthy.

### ClosePositionIfUnderStopLossPrice

The `ClosePositionIfUnderStopLossPrice` function checks if the position lp token price is lower than stopLossPrice, it close the position.

### UpdatePoolHealth

The `UpdatePoolHealth` function updates the health of a pool based on the current state of its positions.

### GetPositionHealth

The `GetPositionHealth` function calculates the health of a position based on its collateral and liabilities.
