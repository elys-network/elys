package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// createElysMultihopExpectedSwapOuts does the same as createMultihopExpectedSwapOuts, however discounts the swap fee
func (k Keeper) createElysMultihopExpectedSwapOuts(
	ctx sdk.Context,
	routes []types.SwapAmountOutRoute,
	tokenOut sdk.Coin,
	cumulativeRouteSwapFee, sumOfSwapFees sdk.Dec,
) ([]sdk.Int, error) {
	insExpected := make([]sdk.Int, len(routes))
	for i := len(routes) - 1; i >= 0; i-- {
		route := routes[i]

		pool, poolExists := k.GetPool(ctx, route.PoolId)
		if !poolExists {
			return nil, types.ErrInvalidPoolId
		}

		swapFee := pool.GetPoolParams().SwapFee
		tokenIn, err := pool.CalcInAmtGivenOut(sdk.NewCoins(tokenOut), route.TokenInDenom, cumulativeRouteSwapFee.Mul((swapFee.Quo(sumOfSwapFees))))
		if err != nil {
			return nil, err
		}

		insExpected[i] = tokenIn.Amount
		tokenOut = tokenIn
	}

	return insExpected, nil
}
