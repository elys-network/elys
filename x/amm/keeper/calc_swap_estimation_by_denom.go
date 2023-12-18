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
	err error,
) {
	// if amount denom is equal to denomIn, calculate swap estimation by denomIn
	if amount.Denom == denomIn {
		inRoute, err := k.CalcInRouteByDenom(ctx, denomIn, denomOut, baseCurrency)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}
		spotPrice, tokenOut, swapFeeOut, _, availableLiquidity, err := k.CalcInRouteSpotPrice(ctx, amount, inRoute, discount, overrideSwapFee)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}
		return inRoute, nil, tokenOut, spotPrice, swapFeeOut, discount, availableLiquidity, nil
	}

	// if amount denom is equal to denomOut, calculate swap estimation by denomOut
	if amount.Denom == denomOut {
		outRoute, err := k.CalcOutRouteByDenom(ctx, denomOut, denomIn, baseCurrency)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}
		spotPrice, tokenIn, swapFeeOut, _, availableLiquidity, err := k.CalcOutRouteSpotPrice(ctx, amount, outRoute, discount, overrideSwapFee)
		if err != nil {
			return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, err
		}
		return nil, outRoute, tokenIn, spotPrice, swapFeeOut, discount, availableLiquidity, nil
	}

	// if amount denom is neither equal to denomIn nor denomOut, return error
	return nil, nil, sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.Coin{}, types.ErrInvalidDenom
}
