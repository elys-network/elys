package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// ClaimReward withdraw specific amount of coin from unclaimed reward
func (k msgServer) ClaimReward(goCtx context.Context, msg *types.MsgClaimReward) (*types.MsgClaimRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Withdraw tokens
	err := k.ProcessClaimReward(ctx, msg.Creator, msg.Denom, msg.Amount)

	return &types.MsgClaimRewardResponse{}, err
}
