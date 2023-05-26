<!--
order: 3
-->

# Keepers

## WeightBreakingFee

```go
WeightDiffAvg = Average(abs(Weight - TargetWeight))
WeightBroken = max(WeightDiffAvg(AfterAction) - WeightDiffAvg(BeforeAction), 0)
WeightBreakingFee = WeightBroken \* WeightBreakingFeeMutliplier
```

## Fee distribution logic

Fees are transferred to fee collector wallet and distributed per epoch to LPs and stakers per epoch.

LPs claim rewards through `commitment` module.
For stakers, will need to distribute through `incentive` module.

## Swap

### Osmosis swap logic

```go
func (server msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	tokenOutAmount, err := server.keeper.MultihopSwapExactAmountIn(ctx, sender, msg.Routes, msg.TokenIn, msg.TokenOutMinAmount)
	if err != nil {
		return nil, err
	}

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgSwapExactAmountInResponse{TokenOutAmount: tokenOutAmount}, nil
}

func (server msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	tokenInAmount, err := server.keeper.MultihopSwapExactAmountOut(ctx, sender, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
	if err != nil {
		return nil, err
	}

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgSwapExactAmountOutResponse{TokenInAmount: tokenInAmount}, nil
}
```

### Swap logic change on Elys

Swap logic should be different per oracle case and no oracle case. Weight is dynamically determined from oracle and based on new weight, swap out amount is determined.

```go
Slippage = (InAmount/Spot_Price - Swap_out_amount) / (InAmount/Spot_Price)
ActualSlippage = Slippage * params.SlippageReduction
```

### Fees on swap operation

```
UseOracle(false): Swap Fee
UseOracle(true): Swap Fee + WeightBreakFee
```

## Add Liquidity

### Osmosis add liquidity logic

```go
func (server msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool) (*types.MsgJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	neededLp, sharesOut, err := server.keeper.JoinPoolNoSwap(ctx, sender, msg.PoolId, msg.ShareOutAmount, msg.TokenInMaxs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgJoinPoolResponse{
		ShareOutAmount: sharesOut,
		TokenIn:        neededLp,
	}, nil
}
```

### Fees on add liquidity operation

```
UseOracle(false): 0
UseOracle(true): WeightBreakFee
```

## Remove Liquidity

### Osmosis remove liquidity logic

```go
func (server msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	exitCoins, err := server.keeper.ExitPool(ctx, sender, msg.PoolId, msg.ShareInAmount, msg.TokenOutMins)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgExitPoolResponse{
		TokenOut: exitCoins,
	}, nil
}
```

### Fees on remove liquidity operation

```
UseOracle(false): Exit Fee
UseOracle(true): Exit Fee + WeightBreakFee
```
