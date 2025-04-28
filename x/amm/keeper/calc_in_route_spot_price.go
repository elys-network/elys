package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInRouteSpotPrice calculates the spot price of the given token and in route
func (k Keeper) CalcInRouteSpotPrice(ctx sdk.Context,
	tokenIn sdk.Coin,
	routes []*types.SwapAmountInRoute,
	discount sdkmath.LegacyDec,
	overrideSwapFee sdkmath.LegacyDec,
) (sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coin, sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
	if len(routes) == 0 {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrEmptyRoutes
	}

	// Start with the initial token input
	tokensIn := sdk.Coins{tokenIn}

	// The final token out denom
	var tokenOutDenom string

	// convert routes []*types.SwapAmountInRoute to []types.SwapAmountInRoute
	routesNoPtr := make([]types.SwapAmountInRoute, len(routes))
	for i, route := range routes {
		routesNoPtr[i] = *route
	}

	route := types.SwapAmountInRoutes(routesNoPtr)
	if err := route.Validate(); err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee := sdkmath.LegacyZeroDec()

	// Track the total available liquidity in the pool for final token out denom
	var availableLiquidity sdk.Coin

	weightBalance := sdkmath.LegacyZeroDec()
	slippage := sdkmath.LegacyZeroDec()

	for _, route := range routes {
		poolId := route.PoolId
		tokenOutDenom = route.TokenOutDenom

		pool, found := k.GetPool(ctx, poolId)
		if !found {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrPoolNotFound
		}

		// Get Pool swap fee
		// Divide fees with the number of routes to incentivize multi-hop
		swapFee := pool.GetPoolParams().SwapFee.Quo(sdkmath.LegacyNewDec(int64(len(routes))))
		takersFee := k.parameterKeeper.GetParams(ctx).TakerFees.Quo(sdkmath.LegacyNewDec(int64(len(routes))))

		// Override swap fee if applicable
		if overrideSwapFee.IsPositive() {
			swapFee = overrideSwapFee
		}

		// Apply discount to swap fee
		swapFee = types.ApplyDiscount(swapFee, discount)

		// Estimate swap
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		tokenOut, swapSlippage, _, weightBalanceBonus, oracleOutAmount, swapFee, err := k.SwapOutAmtGivenIn(cacheCtx, pool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, swapFee, sdkmath.LegacyOneDec(), takersFee)
		if err != nil {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// Check treasury and update weightBalance
		bonusTokenAmount := sdkmath.ZeroInt()
		if weightBalanceBonus.IsPositive() {
			// get treasury balance
			rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
			treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

			bonusTokenAmount = oracleOutAmount.Mul(weightBalanceBonus).TruncateInt()

			if treasuryTokenAmount.LT(bonusTokenAmount) {
				weightBalanceBonus = treasuryTokenAmount.ToLegacyDec().Quo(oracleOutAmount)
				bonusTokenAmount = treasuryTokenAmount
			}
		}

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		if tokenOut.IsZero() {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrAmountTooLow
		}

		if bonusTokenAmount.IsPositive() {
			tokenOut = sdk.NewCoin(tokenOut.Denom, tokenOut.Amount.Add(bonusTokenAmount))
		}

		// Use the current swap result as the input for the next iteration
		tokensIn = sdk.Coins{tokenOut}

		// Get the available liquidity for the final token out denom
		_, poolAsset, err := pool.GetPoolAssetAndIndex(tokenOutDenom)
		if err != nil {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}
		// Use accounted pool balance
		accAmount := k.accountedPoolKeeper.GetAccountedBalance(ctx, pool.PoolId, poolAsset.Token.Denom)
		if accAmount.IsPositive() {
			poolAsset.Token.Amount = accAmount
		}
		availableLiquidity = poolAsset.Token
		weightBalance = weightBalanceBonus
		slippage = slippage.Add(swapSlippage)
	}

	// Ensure tokenIn.Amount is not zero to avoid division by zero
	if tokenIn.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrAmountTooLow
	}

	// Calculate the spot price given the initial token in and the final token out
	impactedPrice := sdkmath.LegacyNewDecFromInt(tokensIn[0].Amount).Quo(sdkmath.LegacyNewDecFromInt(tokenIn.Amount))

	// Calculate spot price with GetTokenARate
	spotPrice := sdkmath.LegacyOneDec()
	tokenInDenom := tokenIn.Denom
	for _, route := range routes {
		poolId := route.PoolId
		tokenOutDenom = route.TokenOutDenom

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
		tokenInDenom = tokenOutDenom
		spotPrice = spotPrice.Mul(rate)
	}

	if !spotPrice.IsPositive() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), types.ErrSpotPriceIsZero
	}

	// Construct the token out coin
	tokenOut := tokensIn[0]

	return spotPrice, impactedPrice, tokenOut, totalDiscountedSwapFee, discount, availableLiquidity, slippage, weightBalance, nil
}
