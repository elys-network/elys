package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) getElysRoutedMultihopTotalSwapFee(ctx sdk.Context, route types.MultihopRoute) (
	totalPathSwapFee sdkmath.LegacyDec, sumOfSwapFees sdkmath.LegacyDec, err error,
) {
	additiveSwapFee := sdkmath.LegacyZeroDec()
	maxSwapFee := sdkmath.LegacyZeroDec()

	for _, poolId := range route.PoolIds() {
		pool, poolExists := k.GetPool(ctx, poolId)
		if !poolExists {
			return sdkmath.LegacyDec{}, sdkmath.LegacyDec{}, types.ErrInvalidPoolId
		}
		swapFee := pool.GetPoolParams().SwapFee
		additiveSwapFee = additiveSwapFee.Add(swapFee)
		maxSwapFee = sdkmath.LegacyMaxDec(maxSwapFee, swapFee)
	}
	averageSwapFee := additiveSwapFee.QuoInt64(2)
	maxSwapFee = sdkmath.LegacyMaxDec(maxSwapFee, averageSwapFee)
	return maxSwapFee, additiveSwapFee, nil
}

func (k Keeper) isElysRoutedMultihop(_ sdk.Context, route types.MultihopRoute, inDenom, outDenom string) (isRouted bool) {
	if route.Length() != 2 {
		return false
	}
	intemediateDenoms := route.IntermediateDenoms()
	if len(intemediateDenoms) != 1 {
		return false
	}
	if inDenom == outDenom {
		return false
	}

	poolIds := route.PoolIds()

	return poolIds[0] != poolIds[1]
}
