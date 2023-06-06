package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// createMultihopExpectedSwapOuts defines the output denom and output amount for the last pool in
// the route of pools the caller is intending to hop through in a fixed-output multihop tx. It estimates the input
// amount for this last pool and then chains that input as the output of the previous pool in the route, repeating
// until the first pool is reached. It returns an array of inputs, each of which correspond to a pool ID in the
// route of pools for the original multihop transaction.
// TODO: test this.
func (k Keeper) createMultihopExpectedSwapOuts(
	ctx sdk.Context,
	routes []types.SwapAmountOutRoute,
	tokenOut sdk.Coin,
) ([]math.Int, error) {
	insExpected := make([]math.Int, len(routes))
	for i := len(routes) - 1; i >= 0; i-- {
		route := routes[i]

		pool, poolExists := k.GetPool(ctx, route.PoolId)
		if !poolExists {
			return nil, types.ErrInvalidPoolId
		}

		tokenIn, err := pool.CalcInAmtGivenOut(sdk.NewCoins(tokenOut), route.TokenInDenom, pool.GetPoolParams().SwapFee)
		if err != nil {
			return nil, err
		}

		insExpected[i] = tokenIn.Amount
		tokenOut = tokenIn
	}

	return insExpected, nil
}
