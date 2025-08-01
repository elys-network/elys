package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/masterchef/types"
	stabletypes "github.com/elys-network/elys/v7/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

	stable_apr := osmomath.ZeroBigDec()
	if req.PoolId >= stabletypes.UsdcPoolId {
		borrowPool, found := k.stableKeeper.GetPool(ctx, req.PoolId)
		if found {
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{PoolId: req.PoolId})
			if err == nil {
				stable_apr = borrowPool.GetBigDecInterestRate().MulDec(res.BorrowRatio)
			}
		}
	}

	return &types.QueryPoolInfoResponse{PoolInfo: poolInfo, StableApr: stable_apr.Dec()}, nil
}

func (k Keeper) ListPoolInfos(goCtx context.Context, req *types.QueryListPoolInfosRequest) (*types.QueryListPoolInfosResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	poolStore := prefix.NewStore(store, types.PoolInfoKeyPrefix)

	var list []types.QueryPoolInfoResponse

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.PoolInfo
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}

		stable_apr := osmomath.ZeroBigDec()
		if pool.PoolId >= stabletypes.UsdcPoolId {
			borrowPool, found := k.stableKeeper.GetPool(ctx, pool.PoolId)
			if found {
				res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{PoolId: pool.PoolId})
				if err == nil {
					stable_apr = borrowPool.GetBigDecInterestRate().MulDec(res.BorrowRatio)
				}
			}
		}

		list = append(list, types.QueryPoolInfoResponse{PoolInfo: pool, StableApr: stable_apr.Dec()})

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListPoolInfosResponse{List: list, Pagination: pageRes}, nil
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

func (k Keeper) TotalPendingRewards(goCtx context.Context, req *types.QueryTotalPendingRewardsRequest) (*types.QueryTotalPendingRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var totalRewards sdk.Coins
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	positionStore := prefix.NewStore(store, types.UserRewardInfoKeyPrefix)

	if req.Pagination == nil {
		req.Pagination = &query.PageRequest{
			Limit: 100000,
		}
	}

	count := uint64(0)

	pageRes, err := query.Paginate(positionStore, req.Pagination, func(key []byte, value []byte) error {
		var reward types.UserRewardInfo
		k.cdc.MustUnmarshal(value, &reward)
		k.AfterWithdraw(ctx, reward.PoolId, sdk.MustAccAddressFromBech32(reward.User), sdkmath.ZeroInt())
		if reward.RewardPending.IsPositive() {
			totalRewards = totalRewards.Add(sdk.NewCoin(reward.RewardDenom, reward.RewardPending.TruncateInt()))
		}
		count++
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryTotalPendingRewardsResponse{
		TotalPendingRewards: totalRewards,
		Count:               count,
		Pagination:          pageRes,
	}, nil
}

func (k Keeper) PendingRewards(goCtx context.Context, req *types.QueryPendingRewardsRequest) (*types.QueryPendingRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.UserRewardInfoKeyPrefix)

	defer iterator.Close()

	var totalRewards sdk.Coins
	count := uint64(0)

	for ; iterator.Valid(); iterator.Next() {
		var val types.UserRewardInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.RewardPending.IsPositive() {
			totalRewards = totalRewards.Add(sdk.NewCoin(val.RewardDenom, val.RewardPending.TruncateInt()))
		}
		count++
	}

	return &types.QueryPendingRewardsResponse{
		TotalPendingRewards: totalRewards,
		Count:               count,
	}, nil
}
