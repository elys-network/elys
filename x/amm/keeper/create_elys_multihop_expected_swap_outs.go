package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// createElysMultihopExpectedSwapOuts does the same as createMultihopExpectedSwapOuts, however discounts the swap fee
func (k Keeper) createElysMultihopExpectedSwapOuts(
	ctx sdk.Context,
	routes []types.SwapAmountOutRoute,
	tokenOut sdk.Coin,
	cumulativeRouteSwapFee, sumOfSwapFees math.LegacyDec,
) ([]math.Int, error) {
	insExpected := make([]math.Int, len(routes))

	for i := len(routes) - 1; i >= 0; i-- {
		route := routes[i]
		pool, poolExists := k.GetPool(ctx, route.PoolId)
		if !poolExists {
			return nil, types.ErrInvalidPoolId
		}

		swapFee := pool.GetPoolParams().SwapFee
		actualSwapFee := math.LegacyZeroDec()
		if sumOfSwapFees.IsPositive() {
			actualSwapFee = cumulativeRouteSwapFee.Mul(swapFee.Quo(sumOfSwapFees))
		}

		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		tokenIn, _, err := pool.CalcInAmtGivenOut(ctx, k.oracleKeeper, &snapshot, sdk.NewCoins(tokenOut), route.TokenInDenom, actualSwapFee, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}

		insExpected[i] = tokenIn.Amount
		tokenOut = tokenIn
	}

	return insExpected, nil
}
