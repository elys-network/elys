package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) getElysRoutedMultihopTotalSwapFee(ctx sdk.Context, route types.MultihopRoute) (
	totalPathSwapFee sdk.Dec, sumOfSwapFees sdk.Dec, err error,
) {
	additiveSwapFee := sdk.ZeroDec()
	maxSwapFee := sdk.ZeroDec()

	for _, poolId := range route.PoolIds() {
		pool, poolExists := k.GetPool(ctx, poolId)
		if !poolExists {
			return sdk.Dec{}, sdk.Dec{}, types.ErrInvalidPoolId
		}
		swapFee := pool.GetPoolParams().SwapFee
		additiveSwapFee = additiveSwapFee.Add(swapFee)
		maxSwapFee = sdk.MaxDec(maxSwapFee, swapFee)
	}
	averageSwapFee := additiveSwapFee.QuoInt64(2)
	maxSwapFee = sdk.MaxDec(maxSwapFee, averageSwapFee)
	return maxSwapFee, additiveSwapFee, nil
}

func (k Keeper) isElysRoutedMultihop(ctx sdk.Context, route types.MultihopRoute, inDenom, outDenom string) (isRouted bool) {
	if route.Length() != 2 {
		return false
	}
	intemediateDenoms := route.IntermediateDenoms()
	if len(intemediateDenoms) != 1 /*|| intemediateDenoms[0] != appparams.BaseCoinUnit*/ {
		return false
	}
	if inDenom == outDenom {
		return false
	}
	poolIds := route.PoolIds()
	if poolIds[0] == poolIds[1] {
		return false
	}

	// route0Incentivized := k.poolIncentivesKeeper.IsPoolIncentivized(ctx, poolIds[0])
	// route1Incentivized := k.poolIncentivesKeeper.IsPoolIncentivized(ctx, poolIds[1])

	// return route0Incentivized && route1Incentivized

	return true
}
