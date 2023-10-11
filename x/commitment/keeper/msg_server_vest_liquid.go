package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) VestLiquid(goCtx context.Context, msg *types.MsgVestLiquid) (*types.MsgVestLiquidResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.DepositLiquidTokensUncommitted(ctx, msg.Denom, msg.Amount, msg.Creator); err != nil {
		return &types.MsgVestLiquidResponse{}, err
	}

	if err := k.ProcessTokenVesting(ctx, msg.Denom, msg.Amount, msg.Creator); err != nil {
		return &types.MsgVestLiquidResponse{}, err
	}

	return &types.MsgVestLiquidResponse{}, nil
}
