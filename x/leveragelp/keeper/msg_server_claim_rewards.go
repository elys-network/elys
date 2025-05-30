package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	for _, id := range msg.Ids {
		position, err := k.GetPosition(ctx, sender, id)
		if err != nil {
			return nil, err
		}

		// Add trigger function
		pool, _ := k.GetPool(ctx, position.AmmPoolId)
		_, closeAttempted, _, err := k.CheckAndLiquidateUnhealthyPosition(ctx, &position, pool)
		if closeAttempted && err != nil {
			return nil, err
		}

		if !(closeAttempted && err == nil) {
			err = k.masterchefKeeper.ClaimRewards(ctx, position.GetPositionAddress(), []uint64{position.AmmPoolId}, sender)
			if err != nil {
				return nil, err
			}
		}
	}

	return &types.MsgClaimRewardsResponse{}, nil
}
