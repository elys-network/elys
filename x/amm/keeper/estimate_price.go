package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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
	snapshot := k.GetPoolWithAccountedBalance(ctx, pool.PoolId)

	rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, baseCurrency, k.accountedPoolKeeper)
	if err != nil {
		return osmomath.ZeroBigDec()
	}

	return rate
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
	tokenIn := sdk.NewCoin(denom, math.NewInt(utils.Pow10(decimal).TruncateInt64()))
	discount := osmomath.OneBigDec()
	spotPrice, _, _, _, _, _, _, _, err := k.CalcInRouteSpotPrice(ctx, tokenIn, routes, discount, osmomath.ZeroBigDec())
	if err != nil {
		return osmomath.ZeroBigDec()
	}
	return spotPrice.Mul(usdcPrice)
}
