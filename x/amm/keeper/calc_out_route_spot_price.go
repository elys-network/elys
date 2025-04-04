package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcOutRouteSpotPrice calculates the spot price of the given token and out route
func (k Keeper) CalcOutRouteSpotPrice(ctx sdk.Context, tokenOut sdk.Coin, routes []*types.SwapAmountOutRoute, discount sdkmath.LegacyDec, overrideSwapFee sdkmath.LegacyDec) (elystypes.Dec34, elystypes.Dec34, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, elystypes.Dec34, elystypes.Dec34, error) {
	if len(routes) == 0 {
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensOut := sdk.Coins{tokenOut}

	// The final token in denom
	var tokenInDenom string

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := sdkmath.LegacyZeroDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	weightBonus := elystypes.ZeroDec34()
	slippage := elystypes.ZeroDec34()

	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), types.ErrPoolNotFound
		}

		// Get Pool swap fee
		swapFee := pool.GetPoolParams().SwapFee.Quo(sdkmath.LegacyNewDec(int64(len(routes))))
		takersFee := k.parameterKeeper.GetParams(ctx).TakerFees.Quo(sdkmath.LegacyNewDec(int64(len(routes))))

		// Override swap fee if applicable
		if overrideSwapFee.IsPositive() {
			swapFee = overrideSwapFee
		}

		// Apply discount
		swapFee = types.ApplyDiscount(swapFee, discount)

		// Estimate swap
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		swapResult, swapSlippage, _, weightBalanceBonus, _, swapFee, err := k.SwapInAmtGivenOut(cacheCtx, pool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, swapFee, sdkmath.LegacyOneDec(), takersFee)
		if err != nil {
			return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		if weightBalanceBonus.IsPositive() {
			rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
			treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

			bonusTokenAmount := elystypes.NewDec34FromInt(tokenOut.Amount).Mul(weightBalanceBonus).ToInt()

			if treasuryTokenAmount.LT(bonusTokenAmount) {
				weightBalanceBonus = elystypes.NewDec34FromInt(treasuryTokenAmount).Quo(elystypes.NewDec34FromInt(tokenOut.Amount))
			}
		}

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		if swapResult.IsZero() {
			return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), types.ErrAmountTooLow
		}

		// Use the current swap result as the input for the next iteration
		tokensOut = sdk.Coins{swapResult}

		// Get the available liquidity for the final token in denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenInDenom)
		if err != nil {
			return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}
		availableLiquidity = poolAsset.Token
		weightBonus = weightBonus.Add(weightBalanceBonus)
		slippage = slippage.Add(swapSlippage)
	}

	// Ensure tokenIn.Amount is not zero to avoid division by zero
	if tokenOut.IsZero() {
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), types.ErrAmountTooLow
	}

	// Calculate the spot price given the initial token in and the final token in
	impactedPrice := elystypes.NewDec34FromInt(tokensOut[0].Amount).QuoInt(tokenOut.Amount)

	// Calculate spot price with GetTokenARate
	spotPrice := elystypes.OneDec34()
	tokenOutDenom := tokenOut.Denom
	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), types.ErrPoolNotFound
		}

		// Estimate swap
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, tokenOutDenom, k.accountedPoolKeeper)
		if err != nil {
			return elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		// set new tokenIn denom for multihop
		tokenOutDenom = tokenInDenom
		spotPrice = spotPrice.Mul(rate)
	}

	// Calculate the token in amount
	tokenIn := tokensOut[0]

	return spotPrice, impactedPrice, tokenIn, totalDiscountedSwapFee, discount, availableLiquidity, slippage, weightBonus, nil
}
