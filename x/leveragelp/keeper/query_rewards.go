package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Rewards(goCtx context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	creator := sdk.MustAccAddressFromBech32(req.Address)

	userPositions := k.GetPositionsForAddress(ctx, creator)

	totalRewards := sdk.Coins{}
	rewardInfos := []*types.RewardInfo{}
	for _, position := range userPositions {

		coins := k.masterchefKeeper.UserPoolPendingReward(ctx, position.GetPositionAddress(), position.AmmPoolId)
		rewardInfos = append(rewardInfos, &types.RewardInfo{
			PoolId:     position.AmmPoolId,
			PositionId: position.Id,
			Reward:     coins,
		})
		totalRewards = totalRewards.Add(coins...)
	}

	return &types.QueryRewardsResponse{Rewards: rewardInfos, TotalRewards: totalRewards}, nil
}
