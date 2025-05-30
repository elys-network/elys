package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Rewards(goCtx context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	totalRewards := sdk.Coins{}
	rewardInfos := []*types.RewardInfo{}
	for _, id := range req.Ids {
		creator := sdk.MustAccAddressFromBech32(req.Address)
		position, err := k.GetPosition(ctx, creator, id)
		if err != nil {
			return &types.QueryRewardsResponse{}, nil
		}
		coins := k.masterchefKeeper.UserPoolPendingReward(ctx, position.GetPositionAddress(), position.AmmPoolId)
		rewardInfos = append(rewardInfos, &types.RewardInfo{
			PositionId: id,
			Reward:     coins,
		})
		totalRewards = totalRewards.Add(coins...)
	}

	return &types.QueryRewardsResponse{Rewards: rewardInfos, TotalRewards: totalRewards}, nil
}
