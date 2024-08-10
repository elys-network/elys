package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	for _, id := range msg.Ids {
		position, err := k.GetPosition(ctx, sender, id)
		if err != nil {
			return nil, err
		}
		posAddr := types.GetPositionAddress(id)
		if position.Address != msg.Sender {
			return nil, types.ErrPositionDoesNotExist
		}
		err = k.masterchefKeeper.ClaimRewards(ctx, posAddr, []uint64{position.AmmPoolId}, sender)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
