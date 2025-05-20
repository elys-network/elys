package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	"github.com/elys-network/elys/v4/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PoolRewards(goCtx context.Context, req *types.QueryPoolRewardsRequest) (*types.QueryPoolRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pools := make([]types.PoolRewards, 0)
	skipCount := uint64(0)

	// calculate external rewards APR
	externalRewardsAprs := k.generateExternalRewardsApr(ctx)

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
			earnPool := k.generatePoolRewards(ctx, &pool, externalRewardsAprs)
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
			poolRewards := k.generatePoolRewards(ctx, &p, externalRewardsAprs)
			pools = append(pools, poolRewards)

			return false
		})
	}

	return &types.QueryPoolRewardsResponse{
		Pools: pools,
	}, nil
}

// Generate earn pool struct
func (k *Keeper) generatePoolRewards(ctx sdk.Context, ammPool *ammtypes.Pool, externalRewardsAprs map[uint64]osmomath.BigDec) types.PoolRewards {
	// Get rewards amount
	rewardsUsd, rewardCoins := k.GetDailyRewardsAmountForPool(ctx, ammPool.PoolId)
	edenForward := sdk.NewCoin(ptypes.Eden, k.ForwardEdenCalc(ctx, ammPool.PoolId).Dec().RoundInt())
	tvl, err := ammPool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	apr := rewardsUsd.MulInt64(365)
	if err != nil {
		apr = osmomath.ZeroBigDec()
	} else {
		apr = apr.Quo(tvl)
	}

	externalRewardAprsLegacyDec := math.LegacyZeroDec()
	if externalRewardApr, ok := externalRewardsAprs[ammPool.PoolId]; ok {
		externalRewardAprsLegacyDec = externalRewardApr.Dec()
	}

	return types.PoolRewards{
		PoolId:             ammPool.PoolId,
		RewardsUsd:         rewardsUsd.Dec(),
		RewardCoins:        rewardCoins,
		EdenForward:        edenForward,
		RewardsUsdApr:      apr.Dec(),
		ExternalRewardsApr: externalRewardAprsLegacyDec,
	}
}

func (k Keeper) generateExternalRewardsApr(ctx sdk.Context) map[uint64]osmomath.BigDec {
	externalIncentives := k.GetAllExternalIncentives(ctx)
	rewardsPerPool := make(map[uint64]osmomath.BigDec)
	curBlockHeight := ctx.BlockHeight()
	totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)

	for _, externalIncentive := range externalIncentives {
		if externalIncentive.FromBlock < curBlockHeight && curBlockHeight <= externalIncentive.ToBlock {
			totalAmount := osmomath.BigDecFromSDKInt(externalIncentive.AmountPerBlock.MulRaw(totalBlocksPerYear))
			price := k.oracleKeeper.GetDenomPrice(ctx, externalIncentive.RewardDenom)

			rewardsPerPool[externalIncentive.PoolId] = rewardsPerPool[externalIncentive.PoolId].Add(totalAmount.Mul(price))
		}
	}

	// Convert all rewards to APR
	// Traverse rewardsPerPool map
	for key, value := range rewardsPerPool {
		// Get total pool liquidity
		pool, found := k.amm.GetPool(ctx, key)
		if !found {
			continue
		}
		totalLiquidity, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			rewardsPerPool[key] = osmomath.ZeroBigDec()
		}
		externalRewardsApr := value.Quo(totalLiquidity)
		rewardsPerPool[key] = externalRewardsApr
	}
	return rewardsPerPool
}
