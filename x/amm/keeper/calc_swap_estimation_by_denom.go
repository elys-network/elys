package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcSwapEstimationByDenom calculates the swap estimation by denom
func (k Keeper) CalcSwapEstimationByDenom(
	ctx sdk.Context,
	amount sdk.Coin,
	denomIn string,
	denomOut string,
	baseCurrency string,
	discount sdk.Dec,
	overrideSwapFee sdk.Dec,
	decimals uint64,
) (
	inRoute []*types.SwapAmountInRoute,
	outRoute []*types.SwapAmountOutRoute,
	outAmount sdk.Coin,
	spotPrice sdk.Dec,
	swapFeeOut sdk.Dec,
	discountOut sdk.Dec,
	availableLiquidity sdk.Coin,
	weightBonus sdk.Dec,
	priceImpact sdk.Dec,
	err error,
) {
	var (
		initialSpotPrice sdk.Dec
	)

	// Initialize return variables
	inRoute, outRoute = nil, nil
	outAmount, availableLiquidity = sdk.Coin{}, sdk.Coin{}
	spotPrice, swapFeeOut, discountOut, weightBonus, priceImpact = sdk.ZeroDec(), sdk.ZeroDec(), discount, sdk.ZeroDec(), sdk.ZeroDec()

	// Determine the correct route based on the amount's denom
	if amount.Denom == denomIn {
		inRoute, err = k.CalcInRouteByDenom(ctx, denomIn, denomOut, baseCurrency)
	} else if amount.Denom == denomOut {
		outRoute, err = k.CalcOutRouteByDenom(ctx, denomOut, denomIn, baseCurrency)
	} else {
		err = types.ErrInvalidDenom
		return
	}

	if err != nil {
		return
	}

	// Calculate initial spot price and price impact if decimals is not zero
	if decimals != 0 {
		lowestAmountForInitialSpotPriceCalc := int64(math.Pow10(int(decimals)))
		initialCoin := sdk.NewInt64Coin(amount.Denom, lowestAmountForInitialSpotPriceCalc)

		if amount.Denom == denomIn {
			initialSpotPrice, _, _, _, _, _, err = k.CalcInRouteSpotPrice(ctx, initialCoin, inRoute, discount, overrideSwapFee)
		} else {
			initialSpotPrice, _, _, _, _, _, err = k.CalcOutRouteSpotPrice(ctx, initialCoin, outRoute, discount, overrideSwapFee)
		}

		if err != nil {
			return
		}
		if initialSpotPrice.IsZero() {
			err = types.ErrInitialSpotPriceIsZero
			return
		}
	}

	// Calculate final spot price and other outputs
	if amount.Denom == denomIn {
		spotPrice, outAmount, swapFeeOut, _, availableLiquidity, weightBonus, err = k.CalcInRouteSpotPrice(ctx, amount, inRoute, discount, overrideSwapFee)
	} else {
		spotPrice, outAmount, swapFeeOut, _, availableLiquidity, weightBonus, err = k.CalcOutRouteSpotPrice(ctx, amount, outRoute, discount, overrideSwapFee)
	}

	if err != nil {
		return
	}

	// Calculate price impact if decimals is not zero
	if decimals != 0 {
		priceImpact = initialSpotPrice.Sub(spotPrice).Quo(initialSpotPrice)
	}

	// Return the calculated values
	return
}
