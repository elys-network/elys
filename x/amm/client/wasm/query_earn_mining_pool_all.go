package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// Check pool type
func (oq *Querier) checkFilterType(ctx sdk.Context, ammPool *types.Pool, filterType types.FilterType) bool {
	switch filterType {
	case types.FilterType_FilterAll:
		return true
	case types.FilterType_FilterPerpetual:
		_, found := oq.perpetualKeeper.GetPool(ctx, ammPool.PoolId)
		return found
	case types.FilterType_FilterFixedWeight:
		return ammPool.PoolParams.UseOracle
	case types.FilterType_FilterDynamicWeight:
		return !ammPool.PoolParams.UseOracle
	case types.FilterType_FilterLeverage:
		_, found := oq.leveragelpKeeper.GetPool(ctx, ammPool.PoolId)
		return found
	}

	return false
}

// Generate earn pool struct
func (oq *Querier) generateEarnPool(ctx sdk.Context, ammPool *types.Pool, filterType types.FilterType) types.EarnPool {
	dexApr := sdk.ZeroDec()
	borrowApr := sdk.ZeroDec()
	leverageLpPercent := sdk.ZeroDec()
	perpetualPercent := sdk.ZeroDec()

	poolInfo, found := oq.incentiveKeeper.GetPoolInfo(ctx, ammPool.PoolId)
	if found {
		dexApr = poolInfo.DexApr
	}

	if filterType == types.FilterType_FilterLeverage {
		prams := oq.stablestakeKeeper.GetParams(ctx)
		borrowApr = prams.InterestRate
	}
	tvl, _ := ammPool.TVL(ctx, oq.oraclekeeper)

	// Get rewards amount
	rewards := oq.incentiveKeeper.GetDexRewardsAmountForPool(ctx, ammPool.PoolId)

	// Get pool ratio
	poolRatio := oq.incentiveKeeper.CalculatePoolRatio(ctx, ammPool)

	leverageLpPool, found := oq.leveragelpKeeper.GetPool(ctx, ammPool.PoolId)
	if found {
		leverageLpPercent = leverageLpPool.Health
	}

	perpetualPool, found := oq.perpetualKeeper.GetPool(ctx, ammPool.PoolId)
	if found {
		perpetualPercent = perpetualPool.Health
	}

	return types.EarnPool{
		Assets:     ammPool.PoolAssets,
		PoolRatio:  poolRatio,
		RewardsApr: dexApr,
		BorrowApr:  borrowApr,
		LeverageLp: leverageLpPercent,
		Perpetual:  perpetualPercent,
		Tvl:        tvl,
		Rewards:    rewards,
		PoolId:     ammPool.PoolId,
		TotalShares: ammPool.TotalShares,
	}
}

// Reverse pools
func (oq *Querier) reversePools(earnPools []types.EarnPool) []types.EarnPool {
	for i, j := 0, len(earnPools)-1; i < j; i, j = i+1, j-1 {
		earnPools[i], earnPools[j] = earnPools[j], earnPools[i]
	}

	return earnPools
}

func (oq *Querier) queryEarnMiningPoolAll(ctx sdk.Context, poolRequest *types.QueryEarnPoolRequest) ([]byte, error) {
	pools := make([]types.EarnPool, 0)
	skipCount := uint64(0)

	// If we have pool Ids specified,
	if len(poolRequest.PoolIds) > 0 {
		for _, pId := range poolRequest.PoolIds {
			pool, found := oq.keeper.GetPool(ctx, pId)
			if !found {
				continue
			}

			// apply filter type
			if !oq.checkFilterType(ctx, &pool, poolRequest.FilterType) {
				continue
			}

			// check offset if pagination available
			if poolRequest.Pagination != nil && skipCount < poolRequest.Pagination.Offset {
				skipCount++
				continue
			}

			// check maximum limit if pagination available
			if poolRequest.Pagination != nil && len(pools) >= int(poolRequest.Pagination.Limit) {
				break
			}

			// Construct earn pool
			earnPool := oq.generateEarnPool(ctx, &pool, poolRequest.FilterType)
			pools = append(pools, earnPool)
		}
	} else {
		oq.keeper.IterateLiquidityPools(ctx, func(p types.Pool) bool {
			// apply filter type
			if !oq.checkFilterType(ctx, &p, poolRequest.FilterType) {
				return false
			}

			// check offset if pagination available
			if poolRequest.Pagination != nil && skipCount < poolRequest.Pagination.Offset {
				skipCount++
				return false
			}

			// check maximum limit if pagination available
			if poolRequest.Pagination != nil && len(pools) >= int(poolRequest.Pagination.Limit) {
				return true
			}

			// Construct earn pool
			earnPool := oq.generateEarnPool(ctx, &p, poolRequest.FilterType)
			pools = append(pools, earnPool)

			return false
		})
	}

	// If reverse is set true.
	if poolRequest.Pagination != nil && poolRequest.Pagination.Reverse {
		pools = oq.reversePools(pools)
	}

	res := types.QueryEarnPoolResponse{
		Pools: pools,
	}

	responseBytes, err := json.Marshal(res)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to serialize pool all response")
	}
	return responseBytes, nil
}
