package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcSwapEstimationByDenom calculates the swap estimation by denom
func (k Keeper) CalcSwapEstimationByDenom(
	ctx sdk.Context,
	amount sdk.Coin,
	denomIn string,
	denomOut string,
	baseCurrency string,
	address string,
	overrideSwapFee sdkmath.LegacyDec,
	decimals uint64,
) (
	inRoute []*types.SwapAmountInRoute,
	outRoute []*types.SwapAmountOutRoute,
	outAmount sdk.Coin,
	spotPrice elystypes.Dec34,
	swapFeeOut sdkmath.LegacyDec,
	discountOut sdkmath.LegacyDec,
	availableLiquidity sdk.Coin,
	slippage elystypes.Dec34,
	weightBonus elystypes.Dec34,
	priceImpact elystypes.Dec34,
	err error,
) {
	var (
		impactedPrice elystypes.Dec34
	)

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)

	// Initialize return variables
	inRoute, outRoute = nil, nil
	outAmount, availableLiquidity = sdk.Coin{}, sdk.Coin{}
	spotPrice, swapFeeOut, discountOut, weightBonus, priceImpact = elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), tier.Discount, elystypes.ZeroDec34(), elystypes.ZeroDec34()

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

	// Calculate final spot price and other outputs
	if amount.Denom == denomIn {
		spotPrice, impactedPrice, outAmount, swapFeeOut, _, availableLiquidity, slippage, weightBonus, err = k.CalcInRouteSpotPrice(ctx, amount, inRoute, tier.Discount, overrideSwapFee)
	} else {
		spotPrice, impactedPrice, outAmount, swapFeeOut, _, availableLiquidity, slippage, weightBonus, err = k.CalcOutRouteSpotPrice(ctx, amount, outRoute, tier.Discount, overrideSwapFee)
	}

	if err != nil {
		return
	}

	// Calculate price impact if decimals is not zero
	if decimals != 0 {
		priceImpact = spotPrice.Sub(impactedPrice).Quo(spotPrice)
	}

	// Return the calculated values
	return
}
