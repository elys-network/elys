package keeper

import (
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
	// if amount denom is equal to denomIn, calculate swap estimation by denomIn
	if amount.Denom == denomIn {
		inRoute, err := k.CalcInRouteByDenom(ctx, denomIn, denomOut, baseCurrency)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		initialSpotPrice, _, _, _, _, _, err := k.CalcInRouteSpotPrice(ctx, sdk.NewCoin(amount.Denom, sdk.ZeroInt()), inRoute, discount, overrideSwapFee)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		// check if initialSpotPrice is zero to avoid division by zero
		if initialSpotPrice.IsZero() {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrInitialSpotPriceIsZero
		}
		spotPrice, tokenOut, swapFeeOut, _, availableLiquidity, weightBonus, err := k.CalcInRouteSpotPrice(ctx, amount, inRoute, discount, overrideSwapFee)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		priceImpact := initialSpotPrice.Sub(spotPrice).Quo(initialSpotPrice)
		return inRoute, nil, tokenOut, spotPrice, swapFeeOut, discount, availableLiquidity, weightBonus, priceImpact, nil
	}

	// if amount denom is equal to denomOut, calculate swap estimation by denomOut
	if amount.Denom == denomOut {
		outRoute, err := k.CalcOutRouteByDenom(ctx, denomOut, denomIn, baseCurrency)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		initialSpotPrice, _, _, _, _, _, err := k.CalcOutRouteSpotPrice(ctx, sdk.NewCoin(amount.Denom, sdk.ZeroInt()), outRoute, discount, overrideSwapFee)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		// check if initialSpotPrice is zero to avoid division by zero
		if initialSpotPrice.IsZero() {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrInitialSpotPriceIsZero
		}
		spotPrice, tokenIn, swapFeeOut, _, availableLiquidity, weightBonus, err := k.CalcOutRouteSpotPrice(ctx, amount, outRoute, discount, overrideSwapFee)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}
		priceImpact := initialSpotPrice.Sub(spotPrice).Quo(initialSpotPrice)
		return nil, outRoute, tokenIn, spotPrice, swapFeeOut, discount, availableLiquidity, weightBonus, priceImpact, nil
	}

	// if amount denom is neither equal to denomIn nor denomOut, return error
	return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), types.ErrInvalidDenom
}
