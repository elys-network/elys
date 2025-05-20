package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/amm/types"
)

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return k.Keeper.SwapExactAmountOut(ctx, msg)
}

func (k Keeper) SwapExactAmountOut(ctx sdk.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
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
	tokenInAmount, swapFee, discount, err := k.RouteExactAmountOut(cacheCtx, sender, recipient, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
	if err != nil {
		return nil, err
	}

	lastSwapIndex := k.GetLastSwapRequestIndex(ctx)
	k.SetSwapExactAmountOutRequests(ctx, msg, lastSwapIndex+1)
	k.SetLastSwapRequestIndex(ctx, lastSwapIndex+1)

	return &types.MsgSwapExactAmountOutResponse{
		TokenInAmount: tokenInAmount,
		SwapFee:       swapFee.Dec(),
		Discount:      discount.Dec(),
		Recipient:     recipient.String(),
	}, nil
}
