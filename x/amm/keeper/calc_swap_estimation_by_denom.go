package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcSwapEstimationByDenom calculates the swap estimation by denom
func (k Keeper) CalcSwapEstimationByDenom(
	ctx sdk.Context,
	amount sdk.Coin,
	denomIn string,
	denomOut string,
	baseAssetsDenoms []string,
	address string,
	overrideSwapFee osmomath.BigDec,
	decimals uint64,
) (
	inRoute []*types.SwapAmountInRoute,
	outRoute []*types.SwapAmountOutRoute,
	outAmount sdk.Coin,
	spotPrice osmomath.BigDec,
	swapFeeOut osmomath.BigDec,
	discountOut osmomath.BigDec,
	availableLiquidity sdk.Coin,
	slippage osmomath.BigDec,
	weightBonus osmomath.BigDec,
	priceImpact osmomath.BigDec,
	err error,
) {
	var (
		impactedPrice osmomath.BigDec
	)

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)

	// Initialize return variables
	inRoute, outRoute = nil, nil
	outAmount, availableLiquidity = sdk.Coin{}, sdk.Coin{}
	spotPrice, swapFeeOut, discountOut, weightBonus, priceImpact = osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.BigDecFromDec(tier.Discount), osmomath.ZeroBigDec(), osmomath.ZeroBigDec()

	// Determine the correct route based on the amount's denom
	if amount.Denom == denomIn {
		for _, baseCurrency := range baseAssetsDenoms {
			inRoute, err = k.CalcInRouteByDenom(ctx, denomIn, denomOut, baseCurrency)
			if err == nil {
				break
			}
		}
	} else if amount.Denom == denomOut {
		for _, baseCurrency := range baseAssetsDenoms {
			outRoute, err = k.CalcOutRouteByDenom(ctx, denomOut, denomIn, baseCurrency)
			if err == nil {
				break
			}
		}
	} else {
		err = types.ErrInvalidDenom
		return
	}

	if err != nil {
		return
	}

	// Calculate final spot price and other outputs
	if amount.Denom == denomIn {
		spotPrice, impactedPrice, outAmount, swapFeeOut, _, availableLiquidity, slippage, weightBonus, err = k.CalcInRouteSpotPrice(ctx, amount, inRoute, osmomath.BigDecFromDec(tier.Discount), overrideSwapFee)
	} else {
		spotPrice, impactedPrice, outAmount, swapFeeOut, _, availableLiquidity, slippage, weightBonus, err = k.CalcOutRouteSpotPrice(ctx, amount, outRoute, osmomath.BigDecFromDec(tier.Discount), overrideSwapFee)
	}

	if err != nil {
		return
	}

	// Calculate price impact if decimals is not zero
	if decimals != 0 {
		if spotPrice.IsZero() {
			err = errors.New("spot price is zero in CalcSwapEstimationByDenom")
			return
		}
		priceImpact = spotPrice.Sub(impactedPrice).Quo(spotPrice)
	}

	// Return the calculated values
	return
}
