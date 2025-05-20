package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/amm/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return k.Keeper.SwapExactAmountIn(ctx, msg)
}

func (k Keeper) SwapExactAmountIn(ctx sdk.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		recipient = sender
	}
	// Try executing the tx on cached context environment, to filter invalid transactions out
	cacheCtx, _ := ctx.CacheContext()
	tokenOutAmount, swapFee, discount, err := k.RouteExactAmountIn(cacheCtx, sender, recipient, msg.Routes, msg.TokenIn, msg.TokenOutMinAmount)
	if err != nil {
		return nil, err
	}

	lastSwapIndex := k.GetLastSwapRequestIndex(ctx)
	k.SetSwapExactAmountInRequests(ctx, msg, lastSwapIndex+1)
	k.SetLastSwapRequestIndex(ctx, lastSwapIndex+1)

	return &types.MsgSwapExactAmountInResponse{
		TokenOutAmount: tokenOutAmount,
		SwapFee:        swapFee.Dec(),
		Discount:       discount.Dec(),
		Recipient:      recipient.String(),
	}, nil
}
