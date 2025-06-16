package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/utils"
	"github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// Estimate the price : eg, 1 Eden -> x usdc
func (k Keeper) EstimatePrice(ctx sdk.Context, tokenInDenom, baseCurrency string) osmomath.BigDec {
	// Find a pool that can convert tokenIn to usdc
	pool, found := k.GetBestPoolWithDenoms(ctx, []string{tokenInDenom, baseCurrency}, false)
	if !found {
		return osmomath.ZeroBigDec()
	}

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, tokenInDenom, baseCurrency)
	if err != nil {
		return osmomath.ZeroBigDec()
	}

	return rate
}

// GetEstimatedTokensPriceFromBestPool returns the estimated price of tokens in USD from the best pool
func (k Keeper) GetEstimatedTokensPriceFromBestPool(ctx sdk.Context, tokenInDenom, tokenOutDenom string) (tokenInPrice, tokenOutPrice math.LegacyDec) {

	baseCurrencyEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec()
	}
	baseCurrencyPrice, found := k.oracleKeeper.GetAssetPrice(ctx, baseCurrencyEntry.DisplayName)
	if !found || baseCurrencyPrice.IsZero() {
		baseCurrencyPrice = math.LegacyOneDec()
	}

	if tokenInDenom == baseCurrencyEntry.Denom {
		tokenInPrice = baseCurrencyPrice
	} else {
		tokenInInfo, found := k.assetProfileKeeper.GetEntryByDenom(ctx, tokenInDenom)
		if !found {
			return math.LegacyZeroDec(), math.LegacyZeroDec()
		}
		tokenInDenomPrice := k.EstimatePrice(ctx, tokenInDenom, baseCurrencyEntry.Denom)
		if tokenInInfo.Decimals >= baseCurrencyEntry.Decimals {
			tokenInPrice = tokenInDenomPrice.Mul(utils.Pow10(tokenInInfo.Decimals - baseCurrencyEntry.Decimals)).MulDec(baseCurrencyPrice).Dec()
		} else {
			tokenInPrice = tokenInDenomPrice.Quo(utils.Pow10(baseCurrencyEntry.Decimals - tokenInInfo.Decimals)).MulDec(baseCurrencyPrice).Dec()
		}
	}

	if tokenOutDenom == baseCurrencyEntry.Denom {
		tokenOutPrice = baseCurrencyPrice
	} else {
		tokenOutInfo, found := k.assetProfileKeeper.GetEntryByDenom(ctx, tokenOutDenom)
		if !found {
			return math.LegacyZeroDec(), math.LegacyZeroDec()
		}
		tokenOutDenomPrice := k.EstimatePrice(ctx, tokenOutDenom, baseCurrencyEntry.Denom)
		if tokenOutInfo.Decimals >= baseCurrencyEntry.Decimals {
			tokenOutPrice = tokenOutDenomPrice.Mul(utils.Pow10(tokenOutInfo.Decimals - baseCurrencyEntry.Decimals)).MulDec(baseCurrencyPrice).Dec()
		} else {
			tokenOutPrice = tokenOutDenomPrice.Quo(utils.Pow10(baseCurrencyEntry.Decimals - tokenOutInfo.Decimals)).MulDec(baseCurrencyPrice).Dec()
		}
	}

	return tokenInPrice, tokenOutPrice
}

func (k Keeper) GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) osmomath.BigDec {
	// Calc ueden / uusdc rate
	edenUsdcRate := k.EstimatePrice(ctx, ptypes.Elys, baseCurrency)
	if edenUsdcRate.IsZero() {
		edenUsdcRate = osmomath.OneBigDec()
	}
	usdcDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal := int64(6)
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = int64(usdcEntry.Decimals)
		}
		usdcDenomPrice = osmomath.NewBigDecWithPrec(1, usdcDecimal)
	}
	return edenUsdcRate.Mul(usdcDenomPrice)
}

func (k Keeper) GetTokenPrice(ctx sdk.Context, tokenInDenom, baseCurrency string) osmomath.BigDec {
	oraclePrice := k.oracleKeeper.GetDenomPrice(ctx, tokenInDenom)
	if !oraclePrice.IsZero() {
		return oraclePrice
	}

	// Calc tokenIn / uusdc rate
	tokenUsdcRate := k.EstimatePrice(ctx, tokenInDenom, baseCurrency)
	usdcDenomPrice := k.oracleKeeper.GetDenomPrice(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal := int64(6)
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = int64(usdcEntry.Decimals)
		}
		usdcDenomPrice = osmomath.NewBigDecWithPrec(1, usdcDecimal)
	}
	return tokenUsdcRate.Mul(usdcDenomPrice)
}

func (k Keeper) CalculateCoinsUSDValue(
	ctx sdk.Context,
	coins sdk.Coins,
) osmomath.BigDec {
	totalValueInUSD := osmomath.ZeroBigDec()
	for _, coin := range coins {
		valueInUSD := k.CalculateUSDValue(ctx, coin.Denom, coin.Amount)
		totalValueInUSD = totalValueInUSD.Add(valueInUSD)
	}

	return totalValueInUSD
}

func (k Keeper) CalculateUSDValue(ctx sdk.Context, denom string, amount math.Int) osmomath.BigDec {
	asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, denom)
	if !found {
		osmomath.ZeroBigDec()
	}
	tokenPrice := k.oracleKeeper.GetDenomPrice(ctx, denom)
	if tokenPrice.Equal(osmomath.ZeroBigDec()) {
		tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
	}
	return osmomath.BigDecFromSDKInt(amount).Mul(tokenPrice)
}

// CalcAmmPrice Panics if decimal is > 18, but we do not support >18 as per AddAssetEntry in AssetProfile
func (k Keeper) CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) osmomath.BigDec {
	usdcDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found || denom == usdcDenom {
		return osmomath.ZeroBigDec()
	}
	usdcPrice := k.oracleKeeper.GetDenomPrice(ctx, usdcDenom)
	resp, err := k.InRouteByDenom(ctx, &types.QueryInRouteByDenomRequest{DenomIn: denom, DenomOut: usdcDenom})
	if err != nil {
		return osmomath.ZeroBigDec()
	}

	routes := resp.InRoute
	tokenIn := sdk.NewCoin(denom, math.NewInt(utils.Pow10Int64(decimal)))
	discount := osmomath.OneBigDec()
	spotPrice, _, _, _, _, _, _, _, err := k.CalcInRouteSpotPrice(ctx, tokenIn, routes, discount, osmomath.ZeroBigDec())
	if err != nil {
		return osmomath.ZeroBigDec()
	}
	return spotPrice.Mul(usdcPrice)
}
