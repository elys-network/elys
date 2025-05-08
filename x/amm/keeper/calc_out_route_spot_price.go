package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcOutRouteSpotPrice calculates the spot price of the given token and out route
func (k Keeper) CalcOutRouteSpotPrice(ctx sdk.Context, tokenOut sdk.Coin, routes []*types.SwapAmountOutRoute, discount osmomath.BigDec, overrideSwapFee osmomath.BigDec) (osmomath.BigDec, osmomath.BigDec, sdk.Coin, osmomath.BigDec, osmomath.BigDec, sdk.Coin, osmomath.BigDec, osmomath.BigDec, error) {
	if len(routes) == 0 {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensOut := sdk.Coins{tokenOut}

	// The final token in denom
	var tokenInDenom string

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := osmomath.ZeroBigDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	weightBonus := osmomath.ZeroBigDec()
	slippage := osmomath.ZeroBigDec()

	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrPoolNotFound
		}

		// Get Pool swap fee
		swapFee := pool.GetPoolParams().GetBigDecSwapFee().Quo(osmomath.NewBigDec(int64(len(routes))))
		takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees().Quo(osmomath.NewBigDec(int64(len(routes))))

		// Override swap fee if applicable
		if overrideSwapFee.IsPositive() {
			swapFee = overrideSwapFee
		}

		// Apply discount
		swapFee = types.ApplyDiscount(swapFee, discount)

		// Estimate swap
		snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)
		cacheCtx, _ := ctx.CacheContext()
		swapResult, swapSlippage, _, weightBalanceBonus, _, swapFee, err := k.SwapInAmtGivenOut(cacheCtx, pool.PoolId, k.oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, osmomath.OneBigDec(), takersFee)
		if err != nil {
			return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		bonusTokenAmount := sdkmath.ZeroInt()
		if weightBalanceBonus.IsPositive() {
			rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
			treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

			bonusTokenAmount := osmomath.BigDecFromSDKInt(tokenOut.Amount).Mul(weightBalanceBonus).Dec().TruncateInt()

			if treasuryTokenAmount.LT(bonusTokenAmount) {
				weightBalanceBonus = osmomath.BigDecFromSDKInt(treasuryTokenAmount).Quo(osmomath.BigDecFromSDKInt(tokenOut.Amount))
				bonusTokenAmount = treasuryTokenAmount
			}
		}

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		if swapResult.IsZero() {
			return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrAmountTooLow
		}

		if bonusTokenAmount.IsPositive() {
			swapResult = sdk.NewCoin(swapResult.Denom, swapResult.Amount.Add(bonusTokenAmount))
		}

		// Use the current swap result as the input for the next iteration
		tokensOut = sdk.Coins{swapResult}

		// Get the available liquidity for the final token in denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenInDenom)
		if err != nil {
			return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		availableLiquidity = poolAsset.Token
		weightBonus = weightBalanceBonus
		slippage = slippage.Add(swapSlippage)
	}

	// Ensure tokenIn.Amount is not zero to avoid division by zero
	if tokenOut.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrAmountTooLow
	}

	// Calculate the spot price given the initial token in and the final token in
	impactedPrice := osmomath.BigDecFromSDKInt(tokensOut[0].Amount).Quo(osmomath.BigDecFromSDKInt(tokenOut.Amount))

	// Calculate spot price with GetTokenARate
	spotPrice := osmomath.OneBigDec()
	tokenOutDenom := tokenOut.Denom
	for _, route := range routes {
		poolId := route.PoolId
		tokenInDenom = route.TokenInDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrPoolNotFound
		}

		// Estimate swap
		rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, tokenInDenom, tokenOutDenom)
		if err != nil {
			return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		// set new tokenIn denom for multihop
		tokenOutDenom = tokenInDenom
		spotPrice = spotPrice.Mul(rate)
	}

	// Calculate the token in amount
	tokenIn := tokensOut[0]

	return spotPrice, impactedPrice, tokenIn, totalDiscountedSwapFee, discount, availableLiquidity, slippage, weightBonus, nil
}
