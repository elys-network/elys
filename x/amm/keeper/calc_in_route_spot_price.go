package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInRouteSpotPrice calculates the spot price of the given token and in route
func (k Keeper) CalcInRouteSpotPrice(ctx sdk.Context, tokenIn sdk.Coin, routes []*types.SwapAmountInRoute, discount sdk.Dec) (sdk.Dec, sdk.Coin, sdk.Dec, sdk.Dec, sdk.Coin, error) {
	if routes == nil || len(routes) == 0 {
		return sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensIn := sdk.Coins{tokenIn}

	// The final token out denom
	var tokenOutDenom string

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := sdk.ZeroDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	for _, route := range routes {
		poolId := route.PoolId
		tokenOutDenom = route.TokenOutDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, types.ErrPoolNotFound
		}

		// Get Pool swap fee
		swapFee := pool.GetPoolParams().SwapFee

		// Apply discount to swap fee if applicable
		swapFee, _, err := k.ApplyDiscount(ctx, swapFee, discount, k.BrokerAddress(ctx))
		if err != nil {
			return sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		// Estimate swap
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		swapResult, err := k.CalcOutAmtGivenIn(ctx, pool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, swapFee)

		if err != nil {
			return sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}

		if swapResult.IsZero() {
			return sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, types.ErrAmountTooLow
		}

		// Use the current swap result as the input for the next iteration
		tokensIn = sdk.Coins{swapResult}

		// Get the available liquidity for the final token out denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenOutDenom)
		if err != nil {
			return sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}
		availableLiquidity = poolAsset.Token
	}

	// Calculate the spot price given the initial token in and the final token out
	spotPrice := sdk.NewDecFromInt(tokensIn[0].Amount).Quo(sdk.NewDecFromInt(tokenIn.Amount))

	// Calculate the token out amount
	tokenOutAmt := sdk.NewDecFromInt(tokenIn.Amount).Mul(spotPrice)

	// Construct the token out coin
	tokenOut := sdk.NewCoin(tokenOutDenom, tokenOutAmt.TruncateInt())

	return spotPrice, tokenOut, totalDiscountedSwapFee, discount, availableLiquidity, nil
}
