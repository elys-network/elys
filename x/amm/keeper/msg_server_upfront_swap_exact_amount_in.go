package keeper

import (
	"context"

	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) UpFrontSwapExactAmountIn(goCtx context.Context, msg *types.MsgUpFrontSwapExactAmountIn) (*types.MsgUpFrontSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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

	tokenOutAmount, swapFee, discount, err := k.RouteExactAmountIn(ctx, sender, sender, msg.Routes, msg.TokenIn, msg.TokenOutMinAmount)
	if err != nil {
		return nil, err
	}
	tokenOutCoin := sdk.Coin{Denom: msg.Routes[len(msg.Routes)-1].TokenOutDenom, Amount: tokenOutAmount}
	types.EmitUpFrontSwapEvent(ctx, sender, msg.TokenIn, tokenOutCoin, swapFee.String())

	return &types.MsgUpFrontSwapExactAmountInResponse{
		TokenOutAmount: tokenOutAmount,
		SwapFee:        swapFee.Dec(),
		Discount:       discount.Dec(),
	}, nil
}
