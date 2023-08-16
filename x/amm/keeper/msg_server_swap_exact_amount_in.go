package keeper

import (
	"context"

	// "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lastSwapIndex := k.GetLastSwapRequestIndex(ctx)
	k.SetSwapExactAmountInRequests(ctx, msg, lastSwapIndex+1)
	k.SetLastSwapRequestIndex(ctx, lastSwapIndex+1)
	return &types.MsgSwapExactAmountInResponse{}, nil

	// sender, err := sdk.AccAddressFromBech32(msg.Sender)
	// if err != nil {
	// 	return nil, err
	// }

	// tokenOutAmount, err := k.RouteExactAmountIn(ctx, sender, msg.Routes, msg.TokenIn, math.Int(msg.TokenOutMinAmount))
	// if err != nil {
	// 	return nil, err
	// }

	// // Swap event is handled elsewhere
	// ctx.EventManager().EmitEvents(sdk.Events{
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	// 		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	// 	),
	// })

	// return &types.MsgSwapExactAmountInResponse{TokenOutAmount: tokenOutAmount}, nil
}
