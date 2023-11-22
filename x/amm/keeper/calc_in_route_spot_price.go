package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInRouteSpotPrice calculates the spot price of the given token and in route
func (k Keeper) CalcInRouteSpotPrice(ctx sdk.Context, tokenIn sdk.Coin, routes []*types.SwapAmountInRoute) (sdk.Dec, sdk.Coin, error) {
	if routes == nil || len(routes) == 0 {
		return sdk.ZeroDec(), sdk.Coin{}, types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensIn := sdk.Coins{tokenIn}

	// The final token out denom
	var tokenOutDenom string

	for _, route := range routes {
		poolId := route.PoolId
		tokenOutDenom = route.TokenOutDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdk.ZeroDec(), sdk.Coin{}, types.ErrPoolNotFound
		}

		// Estimate swap
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		swapResult, err := k.CalcOutAmtGivenIn(ctx, pool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, sdk.ZeroDec())

		if err != nil {
			return sdk.ZeroDec(), sdk.Coin{}, err
		}

		if swapResult.IsZero() {
			return sdk.ZeroDec(), sdk.Coin{}, types.ErrAmountTooLow
		}

		// Use the current swap result as the input for the next iteration
		tokensIn = sdk.Coins{swapResult}
	}

	// Calculate the spot price given the initial token in and the final token out
	spotPrice := sdk.NewDecFromInt(tokensIn[0].Amount).Quo(sdk.NewDecFromInt(tokenIn.Amount))

	// Calculate the token out amount
	tokenOutAmt := sdk.NewDecFromInt(tokenIn.Amount).Mul(spotPrice)

	// Construct the token out coin
	tokenOut := sdk.NewCoin(tokenOutDenom, tokenOutAmt.TruncateInt())

	return spotPrice, tokenOut, nil
}
