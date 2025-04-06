package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
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

	poolInfo, found := k.GetPoolInfo(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id")
	}

	stable_apr := sdkmath.LegacyZeroDec()
	if req.PoolId >= stabletypes.UsdcPoolId {
		borrowPool, found := k.stableKeeper.GetPool(ctx, req.PoolId)
		if found {
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{PoolId: req.PoolId})
			if err == nil {
				stable_apr = borrowPool.InterestRate.Mul(res.BorrowRatio)
			}
		}
	}

	return &types.QueryPoolInfoResponse{PoolInfo: poolInfo, StableApr: stable_apr}, nil
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
	user := sdk.MustAccAddressFromBech32(req.User)
	userRewardInfo, found := k.GetUserRewardInfo(ctx, user, req.PoolId, req.RewardDenom)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid pool id or denom")
	}

	return &types.QueryUserRewardInfoResponse{UserRewardInfo: userRewardInfo}, nil
}

func (k Keeper) UserPoolPendingReward(ctx sdk.Context, user sdk.AccAddress, poolId uint64) sdk.Coins {
	k.AfterWithdraw(ctx, poolId, user, sdkmath.ZeroInt())

	poolRewards := sdk.NewCoins()
	for _, rewardDenom := range k.GetRewardDenoms(ctx, poolId) {
		userRewardInfo, found := k.GetUserRewardInfo(ctx, user, poolId, rewardDenom)
		if found && userRewardInfo.RewardPending.IsPositive() {
			poolRewards = poolRewards.Add(
				sdk.NewCoin(
					rewardDenom,
					userRewardInfo.RewardPending.TruncateInt(),
				),
			)
		}
	}
	return poolRewards
}

func (k Keeper) UserPendingReward(goCtx context.Context, req *types.QueryUserPendingRewardRequest) (*types.QueryUserPendingRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	user, err := sdk.AccAddressFromBech32(req.User)
	if err != nil {
		return nil, err
	}

	totalRewards := sdk.NewCoins()
	rewardsInfos := []*types.RewardInfo{}

	for _, pool := range k.GetAllPoolInfos(ctx) {
		poolRewards := k.UserPoolPendingReward(ctx, user, pool.PoolId)
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

	return &types.QueryStableStakeAprResponse{Apr: apr.Dec()}, nil
}

func (k Keeper) PoolAprs(goCtx context.Context, req *types.QueryPoolAprsRequest) (*types.QueryPoolAprsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	data := k.CalculatePoolAprs(ctx, req.PoolIds)
	return &types.QueryPoolAprsResponse{Data: data}, nil
}
