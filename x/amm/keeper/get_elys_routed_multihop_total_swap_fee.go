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
