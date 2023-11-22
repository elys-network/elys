package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcOutRouteSpotPrice calculates the spot price of the given token and out route
func (k Keeper) CalcOutRouteSpotPrice(ctx sdk.Context, tokenOut sdk.Coin, routes []*types.SwapAmountOutRoute) (sdk.Dec, sdk.Coin, error) {
	if routes == nil || len(routes) == 0 {
		return sdk.ZeroDec(), sdk.Coin{}, types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensOut := sdk.Coins{tokenOut}

	// The final token in denom
	var tokenInDenom string

	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdk.ZeroDec(), sdk.Coin{}, types.ErrPoolNotFound
		}

		// Estimate swap
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		swapResult, err := k.CalcInAmtGivenOut(ctx, pool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, sdk.ZeroDec())

		if err != nil {
			return sdk.ZeroDec(), sdk.Coin{}, err
		}

		if swapResult.IsZero() {
			return sdk.ZeroDec(), sdk.Coin{}, types.ErrAmountTooLow
		}

		// Use the current swap result as the input for the next iteration
		tokensOut = sdk.Coins{swapResult}
	}

	// Calculate the spot price given the initial token in and the final token in
	spotPrice := sdk.NewDecFromInt(tokensOut[0].Amount).Quo(sdk.NewDecFromInt(tokenOut.Amount))

	// Calculate the token in amount
	tokenInAmt := sdk.NewDecFromInt(tokenOut.Amount).Mul(spotPrice)

	// Construct the token out coin
	tokenIn := sdk.NewCoin(tokenInDenom, tokenInAmt.TruncateInt())

	return spotPrice, tokenIn, nil
}
