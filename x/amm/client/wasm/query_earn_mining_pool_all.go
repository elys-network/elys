package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (oq *Querier) checkFilterType(ctx sdk.Context, ammPool *types.Pool, filterType types.FilterType) bool {
	switch filterType {
	case types.FilterType_FilterAll:
		return true
	case types.FilterType_FilterPerpetual:
		_, found := oq.marginKeeper.GetPool(ctx, ammPool.PoolId)
		return found
	case types.FilterType_FilterFixedWeight:
		return ammPool.PoolParams.UseOracle
	case types.FilterType_FilterDynamicWeight:
		return !ammPool.PoolParams.UseOracle
	}

	return false
}

// Generate earn pool struct
func (oq *Querier) generateEarnPool(ctx sdk.Context, ammPool *types.Pool) types.EarnPool {
	dexApr := sdk.ZeroDec()
	edenApr := sdk.ZeroDec()
	poolInfo, found := oq.incentiveKeeper.GetPoolInfo(ctx, ammPool.PoolId)
	if found {
		dexApr = poolInfo.DexApr
		edenApr = poolInfo.EdenApr
	}

	tvl, _ := ammPool.TVL(ctx, oq.oraclekeeper)

	// Get rewards amount
	rewards := oq.incentiveKeeper.GetDexRewardsAmountForPool(ctx, ammPool.PoolId)

	// Get pool share
	poolShare, _ := oq.incentiveKeeper.CalculatePoolShare(ctx, ammPool)

	leverageLpPercent := sdk.ZeroDec()
	perpetualPercent := sdk.ZeroDec()

	return types.EarnPool{
		Assets:     ammPool.PoolAssets,
		PoolRatio:  poolShare,
		DexApr:     dexApr,
		EdenApr:    edenApr,
		LeverageLp: leverageLpPercent,
		Perpetual:  perpetualPercent,
		Tvl:        tvl,
		Rewards:    rewards,
	}
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
			earnPool := oq.generateEarnPool(ctx, &pool)
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
			earnPool := oq.generateEarnPool(ctx, &p)
			pools = append(pools, earnPool)

			return false
		})
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
