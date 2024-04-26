package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Apr(goCtx context.Context, req *types.QueryAprRequest) (*types.QueryAprResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	apr, err := k.CalculateApr(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.QueryAprResponse{Apr: apr}, nil
}

func (k Keeper) Aprs(goCtx context.Context, req *types.QueryAprsRequest) (*types.QueryAprsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	usdcAprUsdc, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_USDC_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprUsdc, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_USDC_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	usdcAprEdenb, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDENB_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprEdenb, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDENB_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	usdcAprEden, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDEN_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprEden, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDEN_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	edenbAprEden, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDEN_PROGRAM, Denom: ptypes.EdenB})
	if err != nil {
		return nil, err
	}

	usdcAprElys, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_ELYS_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprElys, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_ELYS_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	edenbAprElys, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_ELYS_PROGRAM, Denom: ptypes.EdenB})
	if err != nil {
		return nil, err
	}

	return &types.QueryAprsResponse{
		UsdcAprUsdc:  usdcAprUsdc,
		EdenAprUsdc:  edenAprUsdc,
		UsdcAprEdenb: usdcAprEdenb,
		EdenAprEdenb: edenAprEdenb,
		UsdcAprEden:  usdcAprEden,
		EdenAprEden:  edenAprEden,
		EdenbAprEden: edenbAprEden,
		UsdcAprElys:  usdcAprElys,
		EdenAprElys:  edenAprElys,
		EdenbAprElys: edenbAprElys,
	}, nil
}

func (k Keeper) PoolRewards(goCtx context.Context, req *types.QueryPoolRewardsRequest) (*types.QueryPoolRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pools := make([]types.PoolRewards, 0)
	skipCount := uint64(0)

	// If we have pool Ids specified,
	if len(req.PoolIds) > 0 {
		for _, pId := range req.PoolIds {
			pool, found := k.amm.GetPool(ctx, pId)
			if !found {
				continue
			}

			// check offset if pagination available
			if req.Pagination != nil && skipCount < req.Pagination.Offset {
				skipCount++
				continue
			}

			// check maximum limit if pagination available
			if req.Pagination != nil && len(pools) >= int(req.Pagination.Limit) {
				break
			}

			// Construct earn pool
			earnPool := k.generatePoolRewards(ctx, &pool)
			pools = append(pools, earnPool)
		}
	} else {
		k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
			// check offset if pagination available
			if req.Pagination != nil && skipCount < req.Pagination.Offset {
				skipCount++
				return false
			}

			// check maximum limit if pagination available
			if req.Pagination != nil && len(pools) >= int(req.Pagination.Limit) {
				return true
			}

			// Construct earn pool
			poolRewards := k.generatePoolRewards(ctx, &p)
			pools = append(pools, poolRewards)

			return false
		})
	}

	return &types.QueryPoolRewardsResponse{
		Pools: pools,
	}, nil
}

func (k Keeper) AllProgramRewards(goCtx context.Context, req *types.QueryAllProgramRewardsRequest) (*types.QueryAllProgramRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.masterchef.AfterWithdraw(ctx, stablestaketypes.PoolId, req.Address, sdk.ZeroInt())

	stableStakeRewards := sdk.NewDecCoins()
	for _, rewardDenom := range k.masterchef.GetRewardDenoms(ctx, stablestaketypes.PoolId) {
		userRewardInfo, found := k.masterchef.GetUserRewardInfo(ctx, req.Address, stablestaketypes.PoolId, rewardDenom)
		if found && userRewardInfo.RewardPending.IsPositive() {
			stableStakeRewards = stableStakeRewards.Add(
				sdk.NewDecCoinFromDec(
					rewardDenom,
					userRewardInfo.RewardPending,
				),
			)
		}
	}

	delAddr := sdk.MustAccAddressFromBech32(req.Address)
	delegations := k.estaking.Keeper.GetDelegatorDelegations(ctx, delAddr, 5000)
	elysStakingRewards := sdk.DecCoins{}
	for _, del := range delegations {
		rewards, err := k.estaking.DelegationRewards(ctx, req.Address, del.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		elysStakingRewards = elysStakingRewards.Add(rewards...)
	}

	// Eden commit rewards
	edenVal := k.estaking.GetParams(ctx).EdenCommitVal
	edenCommitRewards, err := k.estaking.DelegationRewards(ctx, req.Address, edenVal)
	if err != nil {
		return nil, err
	}

	// EdenB commit rewards
	edenBVal := k.estaking.GetParams(ctx).EdenbCommitVal
	edenBCommitRewards, err := k.estaking.DelegationRewards(ctx, req.Address, edenBVal)
	if err != nil {
		return nil, err
	}

	return &types.QueryAllProgramRewardsResponse{
		UsdcStakingRewards:  stableStakeRewards,
		ElysStakingRewards:  elysStakingRewards,
		EdenStakingRewards:  edenCommitRewards,
		EdenbStakingRewards: edenBCommitRewards,
	}, nil
}

// Generate earn pool struct
func (k *Keeper) generatePoolRewards(ctx sdk.Context, ammPool *ammtypes.Pool) types.PoolRewards {
	// Get rewards amount
	rewardsUsd, rewardCoins := k.GetDailyRewardsAmountForPool(ctx, ammPool.PoolId)

	return types.PoolRewards{
		PoolId:      ammPool.PoolId,
		RewardsUsd:  rewardsUsd,
		RewardCoins: rewardCoins,
	}
}
