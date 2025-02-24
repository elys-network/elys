package keeper

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Estimate the price : eg, 1 Eden -> x usdc
func (k Keeper) EstimatePrice(ctx sdk.Context, tokenInDenom, baseCurrency string) math.LegacyDec {
	// Find a pool that can convert tokenIn to usdc
	pool, found := k.GetBestPoolWithDenoms(ctx, []string{tokenInDenom, baseCurrency}, false)
	if !found {
		return math.LegacyZeroDec()
	}

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)

	rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, baseCurrency, k.accountedPoolKeeper)
	if err != nil {
		return math.LegacyZeroDec()
	}

	return rate
}

func (k Keeper) GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) math.LegacyDec {
	// Calc ueden / uusdc rate
	edenUsdcRate := k.EstimatePrice(ctx, ptypes.Elys, baseCurrency)
	if edenUsdcRate.IsZero() {
		edenUsdcRate = math.LegacyOneDec()
	}
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal := int64(6)
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = int64(usdcEntry.Decimals)
		}
		usdcDenomPrice = math.LegacyNewDecWithPrec(1, usdcDecimal)
	}
	return edenUsdcRate.Mul(usdcDenomPrice)
}

func (k Keeper) GetTokenPrice(ctx sdk.Context, tokenInDenom, baseCurrency string) math.LegacyDec {
	oraclePrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if !oraclePrice.IsZero() {
		return oraclePrice
	}

	// Calc tokenIn / uusdc rate
	tokenUsdcRate := k.EstimatePrice(ctx, tokenInDenom, baseCurrency)
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal := int64(6)
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = int64(usdcEntry.Decimals)
		}
		usdcDenomPrice = math.LegacyNewDecWithPrec(1, usdcDecimal)
	}
	return tokenUsdcRate.Mul(usdcDenomPrice)
}

func (k Keeper) CalculateCoinsUsdcValue(
	ctx sdk.Context,
	coins sdk.Coins,
) sdkmath.LegacyDec {
	totalValueInUSD := sdkmath.LegacyZeroDec()
	for _, coin := range coins {
		valueInUSD := k.CalculateUSDValue(ctx, coin.Denom, coin.Amount)
		totalValueInUSD = totalValueInUSD.Add(valueInUSD)
	}

	return totalValueInUSD
}

func (k Keeper) CalculateUSDValue(ctx sdk.Context, denom string, amount sdkmath.Int) sdkmath.LegacyDec {
	asset, found := k.assetProfileKeeper.GetEntryByDenom(ctx, denom)
	if !found {
		sdkmath.LegacyZeroDec()
	}
	tokenPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, denom)
	if tokenPrice.Equal(sdkmath.LegacyZeroDec()) {
		tokenPrice = k.CalcAmmPrice(ctx, asset.Denom, asset.Decimals)
	}
	return amount.ToLegacyDec().Mul(tokenPrice)
}

func (k Keeper) CalcAmmPrice(ctx sdk.Context, denom string, decimal uint64) sdkmath.LegacyDec {
	usdcDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found || denom == usdcDenom {
		return sdkmath.LegacyZeroDec()
	}
	usdcPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, usdcDenom)
	resp, err := k.InRouteByDenom(ctx, &types.QueryInRouteByDenomRequest{DenomIn: denom, DenomOut: usdcDenom})
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}

	routes := resp.InRoute
	tokenIn := sdk.NewCoin(denom, sdkmath.NewInt(Pow10(decimal).TruncateInt64()))
	discount := sdkmath.LegacyNewDec(1)
	spotPrice, _, _, _, _, _, _, _, err := k.CalcInRouteSpotPrice(ctx, tokenIn, routes, discount, sdkmath.LegacyZeroDec())
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}
	return spotPrice.Mul(usdcPrice)
}

func Pow10(decimal uint64) (value sdkmath.LegacyDec) {
	value = sdkmath.LegacyNewDec(1)
	for i := 0; i < int(decimal); i++ {
		value = value.Mul(sdkmath.LegacyNewDec(10))
	}
	return
}
