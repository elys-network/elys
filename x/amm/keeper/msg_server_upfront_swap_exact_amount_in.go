package keeper

import (
	"context"

	"slices"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) UpFrontSwapExactAmountIn(goCtx context.Context, msg *types.MsgUpFrontSwapExactAmountIn) (*types.MsgUpFrontSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return k.Keeper.UpFrontSwapExactAmountIn(ctx, msg)
}

func (k Keeper) UpFrontSwapExactAmountIn(ctx sdk.Context, msg *types.MsgUpFrontSwapExactAmountIn) (*types.MsgUpFrontSwapExactAmountInResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)
	allowed := slices.Contains(params.AllowedUpfrontSwapMakers, msg.Sender)
	if !allowed {
		return nil, types.ErrUnauthorizedUpFrontSwap
	}

	tokenOutAmount, swapFee, discount, err := k.RouteExactAmountIn(ctx, sender, sender, msg.Routes, msg.TokenIn, sdkmath.Int(msg.TokenOutMinAmount))
	if err != nil {
		return nil, err
	}

	return &types.MsgUpFrontSwapExactAmountInResponse{
		TokenOutAmount: tokenOutAmount,
		SwapFee:        swapFee,
		Discount:       discount,
	}, nil
}
