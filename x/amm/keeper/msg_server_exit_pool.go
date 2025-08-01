package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (k msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	exitCoins, weightBalanceBonus, swapFee, slippage, _, err := k.Keeper.ExitPool(ctx, sender, msg.PoolId, msg.ShareAmountIn, msg.MinAmountsOut, msg.TokenOutDenom, false, true)
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
		TokenOut:           exitCoins,
		WeightBalanceRatio: weightBalanceBonus.Dec(),
		SwapFee:            swapFee.Dec(),
		Slippage:           slippage.Dec(),
	}, nil
}
