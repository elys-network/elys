<!--
order: 5
-->

# Functions

## BeginBlocker

The `BeginBlocker` function is called at the beginning of each block to perform necessary updates and maintenance for the `leveragelp` module. It processes health checks and liquidates unhealthy positions if needed.

```go
func (k Keeper) BeginBlocker(ctx sdk.Context) {
    // logic for regular health checks and liquidation
}
```

### LiquidatePositionIfUnhealthy

The `LiquidatePositionIfUnhealthy` function checks the health of a position and liquidates it if it is unhealthy.

```go
func (k Keeper) LiquidatePositionIfUnhealthy(ctx sdk.Context, position *types.Position, pool types.Pool, ammPool ammtypes.Pool) {
    // logic to liquidate unhealthy positions
}
```

### Open

The `Open` function allows a user to open a leveraged position by providing collateral and specifying leverage and a stop-loss price.

```go
func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
    // logic to open a leveraged position
}
```

### Close

The `Close` function allows a user to close a leveraged position, either partially or fully, by specifying the amount of LP tokens to withdraw.

```go
func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
    // logic to close a leveraged position
}
```

### UpdatePoolHealth

The `UpdatePoolHealth` function updates the health of a pool based on the current state of its positions.

```go
func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) {
    // logic to update pool health
}
```

### GetPositionHealth

The `GetPositionHealth` function calculates the health of a position based on its collateral and liabilities.

```go
func (k Keeper) GetPositionHealth(ctx sdk.Context, position types.Position, ammPool ammtypes.Pool) (sdk.Dec, error) {
    // logic to calculate position health
}
```
