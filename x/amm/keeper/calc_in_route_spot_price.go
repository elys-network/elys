package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInRouteSpotPrice calculates the spot price of the given token and in route
func (k Keeper) CalcInRouteSpotPrice(ctx sdk.Context,
	tokenIn sdk.Coin,
	routes []*types.SwapAmountInRoute,
	discount sdk.Dec,
	overrideSwapFee sdk.Dec,
) (sdk.Dec, sdk.Dec, sdk.Coin, sdk.Dec, sdk.Dec, sdk.Coin, sdk.Dec, sdk.Dec, error) {
	if len(routes) == 0 {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensIn := sdk.Coins{tokenIn}

	// The final token out denom
	var tokenOutDenom string

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := sdk.ZeroDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	weightBalance := sdk.ZeroDec()
	slippage := sdk.ZeroDec()

	for _, route := range routes {
		poolId := route.PoolId
		tokenOutDenom = route.TokenOutDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrPoolNotFound
		}

		// Get Pool swap fee
		swapFee := pool.GetPoolParams().SwapFee

		// Override swap fee if applicable
		if overrideSwapFee.IsPositive() {
			swapFee = overrideSwapFee
		}

		// Apply discount to swap fee
		swapFee = types.ApplyDiscount(swapFee, discount)

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		// Estimate swap
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		tokenOut, swapSlippage, _, weightBalanceBonus, err := k.SwapOutAmtGivenIn(cacheCtx, pool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, swapFee)
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		if tokenOut.IsZero() {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrAmountTooLow
		}

		// Use the current swap result as the input for the next iteration
		tokensIn = sdk.Coins{tokenOut}

		// Get the available liquidity for the final token out denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenOutDenom)
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		availableLiquidity = poolAsset.Token
		weightBalance = weightBalance.Add(weightBalanceBonus)
		slippage = slippage.Add(swapSlippage)
	}

	// Ensure tokenIn.Amount is not zero to avoid division by zero
	if tokenIn.IsZero() {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrAmountTooLow
	}

	// Calculate the spot price given the initial token in and the final token out
	impactedPrice := sdk.NewDecFromInt(tokensIn[0].Amount).Quo(sdk.NewDecFromInt(tokenIn.Amount))

	// Calculate spot price with GetTokenARate
	spotPrice := sdk.OneDec()
	tokenInDenom := tokenIn.Denom
	for _, route := range routes {
		poolId := route.PoolId
		tokenOutDenom = route.TokenOutDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrPoolNotFound
		}

		// Estimate swap
		snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
		rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, tokenOutDenom, k.accountedPoolKeeper)
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		// set new tokenIn denom for multihop
		tokenInDenom = tokenOutDenom
		spotPrice = spotPrice.Mul(rate)
	}

	if !spotPrice.IsPositive() {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrSpotPriceIsZero
	}

	// Construct the token out coin
	tokenOut := tokensIn[0]

	return spotPrice, impactedPrice, tokenOut, totalDiscountedSwapFee, discount, availableLiquidity, slippage, weightBalance, nil
}
