package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcOutRouteSpotPrice calculates the spot price of the given token and out route
func (k Keeper) CalcOutRouteSpotPrice(ctx sdk.Context, tokenOut sdk.Coin, routes []*types.SwapAmountOutRoute, discount sdk.Dec, overrideSwapFee sdk.Dec) (sdk.Dec, sdk.Dec, sdk.Coin, sdk.Dec, sdk.Dec, sdk.Coin, sdk.Dec, sdk.Dec, error) {
	if len(routes) == 0 {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensOut := sdk.Coins{tokenOut}

	// The final token in denom
	var tokenInDenom string

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := sdk.ZeroDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	weightBonus := sdk.ZeroDec()
	slippage := sdk.ZeroDec()

	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

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

		// Apply discount
		swapFee = types.ApplyDiscount(swapFee, discount)

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		// Estimate swap
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		swapResult, swapSlippage, _, weightBalanceBonus, err := k.SwapInAmtGivenOut(cacheCtx, pool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, swapFee, math.LegacyOneDec())
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		if swapResult.IsZero() {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrAmountTooLow
		}

		// Use the current swap result as the input for the next iteration
		tokensOut = sdk.Coins{swapResult}

		// Get the available liquidity for the final token in denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenInDenom)
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		availableLiquidity = poolAsset.Token
		weightBonus = weightBonus.Add(weightBalanceBonus)
		slippage = slippage.Add(swapSlippage)
	}

	// Ensure tokenIn.Amount is not zero to avoid division by zero
	if tokenOut.IsZero() {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrAmountTooLow
	}

	// Calculate the spot price given the initial token in and the final token in
	impactedPrice := sdk.NewDecFromInt(tokensOut[0].Amount).Quo(sdk.NewDecFromInt(tokenOut.Amount))

	// Calculate spot price with GetTokenARate
	spotPrice := sdk.OneDec()
	tokenOutDenom := tokenOut.Denom
	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrPoolNotFound
		}

		// Estimate swap
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, tokenOutDenom, k.accountedPoolKeeper)
		if err != nil {
			return sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		// set new tokenIn denom for multihop
		tokenOutDenom = tokenInDenom
		spotPrice = spotPrice.Mul(rate)
	}

	// Calculate the token in amount
	tokenIn := tokensOut[0]

	return spotPrice, impactedPrice, tokenIn, totalDiscountedSwapFee, discount, availableLiquidity, slippage, weightBonus, nil
}
