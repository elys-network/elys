package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
)

// Vest converts user's commitment to vesting - start with unclaimed rewards and if it's not enough deduct from committed bucket
// mainly utilized for Eden
func (k msgServer) Vest(goCtx context.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.ProcessTokenVesting(ctx, msg.Denom, msg.Amount, msg.Creator); err != nil {
		return &types.MsgVestResponse{}, err
	}

	return &types.MsgVestResponse{}, nil
}
