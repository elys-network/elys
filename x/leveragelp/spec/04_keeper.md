<!--
order: 4
-->

# Keeper

## Position Management

The `leveragelp` module's keeper handles the opening, updating, and closing of leveraged positions. It ensures that positions are correctly managed, leverage limits are respected, and health checks are performed to avoid liquidation risks.

### Opening a Position

The `Open` function allows a user to open a leveraged position by providing collateral and specifying leverage and a stop-loss price.

```go
func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
    // logic to open a leveraged position
}
```

### Closing a Position

The `Close` function allows a user to close a leveraged position, either partially or fully, by specifying the amount of LP tokens to withdraw.

```go
func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
    // logic to close a leveraged position
}
```

### Liquidation and Health Checks

The `BeginBlocker` function performs regular health checks on all positions and liquidates unhealthy positions if necessary.

```go
func (k Keeper) BeginBlocker(ctx sdk.Context) {
    // logic for regular health checks and liquidation
}
```
