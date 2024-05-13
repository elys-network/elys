package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ExternalIncentive(goCtx context.Context, req *types.QueryExternalIncentiveRequest) (*types.QueryExternalIncentiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	externalIncentive, found := k.GetExternalIncentive(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid external incentive id")
	}

	return &types.QueryExternalIncentiveResponse{ExternalIncentive: externalIncentive}, nil
}

func (k Keeper) PoolInfo(goCtx context.Context, req *types.QueryPoolInfoRequest) (*types.QueryPoolInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolInfo, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	return &types.QueryPoolInfoResponse{PoolInfo: poolInfo}, nil
}

func (k Keeper) PoolRewardInfo(goCtx context.Context, req *types.QueryPoolRewardInfoRequest) (*types.QueryPoolRewardInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolRewardInfo, found := k.GetPoolRewardInfo(ctx, req.PoolId, req.RewardDenom)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	return &types.QueryPoolRewardInfoResponse{PoolRewardInfo: poolRewardInfo}, nil
}

func (k Keeper) UserRewardInfo(goCtx context.Context, req *types.QueryUserRewardInfoRequest) (*types.QueryUserRewardInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userRewardInfo, found := k.GetUserRewardInfo(ctx, req.User, req.PoolId, req.RewardDenom)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id or denom")
	}

	return &types.QueryUserRewardInfoResponse{UserRewardInfo: userRewardInfo}, nil
}

func (k Keeper) UserPendingReward(goCtx context.Context, req *types.QueryUserPendingRewardRequest) (*types.QueryUserPendingRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	totalRewards := sdk.NewCoins()
	rewardsInfos := []*types.RewardInfo{}

	for _, pool := range k.GetAllPools(ctx) {
		k.AfterWithdraw(ctx, pool.PoolId, req.User, sdk.ZeroInt())

		poolRewards := sdk.NewCoins()
		for _, rewardDenom := range k.GetRewardDenoms(ctx, pool.PoolId) {
			userRewardInfo, found := k.GetUserRewardInfo(ctx, req.User, pool.PoolId, rewardDenom)
			if found && userRewardInfo.RewardPending.IsPositive() {
				poolRewards = poolRewards.Add(
					sdk.NewCoin(
						rewardDenom,
						userRewardInfo.RewardPending.TruncateInt(),
					),
				)
			}
		}
		rewardsInfos = append(rewardsInfos,
			&types.RewardInfo{
				PoolId: pool.PoolId,
				Reward: poolRewards,
			},
		)

		totalRewards = totalRewards.Add(poolRewards...)
	}

	return &types.QueryUserPendingRewardResponse{
		Rewards:      rewardsInfos,
		TotalRewards: totalRewards,
	}, nil
}

func (k Keeper) StableStakeApr(goCtx context.Context, req *types.QueryStableStakeAprRequest) (*types.QueryStableStakeAprResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	apr, err := k.CalculateStableStakeApr(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.QueryStableStakeAprResponse{Apr: apr}, nil
}

func (k Keeper) PoolAprs(goCtx context.Context, req *types.QueryPoolAprsRequest) (*types.QueryPoolAprsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	data := k.CalculatePoolAprs(ctx, req.PoolIds)
	return &types.QueryPoolAprsResponse{Data: data}, nil
}
