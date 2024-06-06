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

## BeginBlocker

The `BeginBlocker` function is called at the beginning of each block to perform necessary updates and maintenance for the `leveragelp` module. It processes health checks and liquidates unhealthy positions if needed.

### LiquidatePositionIfUnhealthy

The `LiquidatePositionIfUnhealthy` function checks the health of a position and liquidates it if it is unhealthy.

### UpdatePoolHealth

The `UpdatePoolHealth` function updates the health of a pool based on the current state of its positions.

### GetPositionHealth

The `GetPositionHealth` function calculates the health of a position based on its collateral and liabilities.
