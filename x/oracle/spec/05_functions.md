<!--
order: 5
-->

# Functions

## EndBlocker

The `EndBlocker` function is called at the end of each block to perform necessary updates and maintenance for the `oracle` module. It processes the removal of outdated prices.

```go
func (k Keeper) EndBlock(ctx sdk.Context) {
    params := k.GetParams(ctx)
    for _, price := range k.GetAllPrice(ctx) {
        if price.Timestamp + params.PriceExpiryTime < uint64(ctx.BlockTime().Unix()) {
            k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
        }
        if price.BlockHeight + params.LifeTimeInBlocks < uint64(ctx.BlockHeight()) {
            k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
        }
    }
}
```

### FeedPrice

The `FeedPrice` function is responsible for submitting a price for an asset.

```go
func (k msgServer) FeedPrice(goCtx context.Context, msg *types.MsgFeedPrice) (*types.MsgFeedPriceResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    feeder, found := k.Keeper.GetPriceFeeder(ctx, msg.Provider)
    if !found {
        return nil, types.ErrNotAPriceFeeder
    }

    if !feeder.IsActive {
        return nil, types.ErrPriceFeederNotActive
    }

    price := types.Price{
        Provider:    msg.Provider,
        Asset:       msg.Asset,
        Price:       msg.Price,
        Source:      msg.Source,
        Timestamp:   uint64(ctx.BlockTime().Unix()),
        BlockHeight: uint64(ctx.BlockHeight()),
    }

    k.SetPrice(ctx, price)
    return &types.MsgFeedPriceResponse{}, nil
}
```

### SetPriceFeeder

The `SetPriceFeeder` function sets the status of a price feeder.

```go
func (k msgServer) SetPriceFeeder(goCtx context.Context, msg *types.MsgSetPriceFeeder) (*types.MsgSetPriceFeederResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    _, found := k.Keeper.GetPriceFeeder(ctx, msg.Feeder)
    if !found {
        return nil, types.ErrNotAPriceFeeder
    }
    k.Keeper.SetPriceFeeder(ctx, types.PriceFeeder{
        Feeder:   msg.Feeder,
        IsActive: msg.IsActive,
    })
    return &types.MsgSetPriceFeederResponse{}, nil
}
```

### RemovePriceFeeder

The `RemovePriceFeeder` function removes a price feeder.

```go
func (k msgServer) DeletePriceFeeder(goCtx context.Context, msg *types.MsgDeletePriceFeeder) (*types.MsgDeletePriceFeederResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    _, found := k.Keeper.GetPriceFeeder(ctx, msg.Feeder)
    if !found {
        return nil, types.ErrNotAPriceFeeder
    }
    k.RemovePriceFeeder(ctx, msg.Feeder)
    return &types.MsgDeletePriceFeederResponse{}, nil
}
```

### FeedMultiplePrices

The `FeedMultiplePrices` function allows a price feeder to submit multiple prices at once.

```go
func (k msgServer) FeedMultiplePrices(goCtx context.Context, msg *types.MsgFeedMultiplePrices) (*types.MsgFeedMultiplePricesResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    feeder, found := k.Keeper.GetPriceFeeder(ctx, msg.Creator)
    if !found {
        return nil, types.ErrNotAPriceFeeder
    }

    if !feeder.IsActive {
        return nil, types.ErrPriceFeederNotActive
    }

    for _, price := range msg.Prices {
        price.Provider = msg.Creator
        price.Timestamp = uint64(ctx.BlockTime().Unix())
        price.BlockHeight = uint64(ctx.BlockHeight())
        k.SetPrice(ctx, price)
    }

    return &types.MsgFeedMultiplePricesResponse{}, nil
}
```
