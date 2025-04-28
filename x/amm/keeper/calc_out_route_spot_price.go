package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcOutRouteSpotPrice calculates the spot price of the given token and out route
func (k Keeper) CalcOutRouteSpotPrice(ctx sdk.Context, tokenOut sdk.Coin, routes []*types.SwapAmountOutRoute, discount sdkmath.LegacyDec, overrideSwapFee sdkmath.LegacyDec) (sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
	if len(routes) == 0 {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensOut := sdk.Coins{tokenOut}

	// The final token in denom
	var tokenInDenom string

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := sdkmath.LegacyZeroDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	weightBonus := sdkmath.LegacyZeroDec()
	slippage := sdkmath.LegacyZeroDec()

	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrPoolNotFound
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
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		bonusTokenAmount := sdkmath.ZeroInt()
		if weightBalanceBonus.IsPositive() {
			rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
			treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

			bonusTokenAmount = tokenOut.Amount.ToLegacyDec().Mul(weightBalanceBonus).TruncateInt()

			if treasuryTokenAmount.LT(bonusTokenAmount) {
				weightBalanceBonus = treasuryTokenAmount.ToLegacyDec().Quo(tokenOut.Amount.ToLegacyDec())
				bonusTokenAmount = treasuryTokenAmount
			}
		}

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		if swapResult.IsZero() {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrAmountTooLow
		}

		if bonusTokenAmount.IsPositive() {
			swapResult = sdk.NewCoin(swapResult.Denom, swapResult.Amount.Add(bonusTokenAmount))
		}

		// Use the current swap result as the input for the next iteration
		tokensOut = sdk.Coins{swapResult}

		// Get the available liquidity for the final token in denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenInDenom)
		if err != nil {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}
		availableLiquidity = poolAsset.Token
		weightBonus = weightBalanceBonus
		slippage = slippage.Add(swapSlippage)
	}

	// Ensure tokenIn.Amount is not zero to avoid division by zero
	if tokenOut.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrAmountTooLow
	}

	// Calculate the spot price given the initial token in and the final token in
	impactedPrice := sdkmath.LegacyNewDecFromInt(tokensOut[0].Amount).Quo(sdkmath.LegacyNewDecFromInt(tokenOut.Amount))

	// Calculate spot price with GetTokenARate
	spotPrice := sdkmath.LegacyOneDec()
	tokenOutDenom := tokenOut.Denom
	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrPoolNotFound
		}

		// Estimate swap
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, tokenOutDenom, k.accountedPoolKeeper)
		if err != nil {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// set new tokenIn denom for multihop
		tokenOutDenom = tokenInDenom
		spotPrice = spotPrice.Mul(rate)
	}

	// Calculate the token in amount
	tokenIn := tokensOut[0]

	return spotPrice, impactedPrice, tokenIn, totalDiscountedSwapFee, discount, availableLiquidity, slippage, weightBonus, nil
}
