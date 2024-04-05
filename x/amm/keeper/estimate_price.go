package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Estimate the price : eg, 1 Eden -> x usdc
func (k Keeper) EstimatePrice(ctx sdk.Context, tokenInDenom, baseCurrency string) math.LegacyDec {
	// Find a pool that can convert tokenIn to usdc
	pool, found := k.GetBestPoolWithDenoms(ctx, []string{tokenInDenom, baseCurrency})
	if !found {
		return sdk.ZeroDec()
	}

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	snapshot := k.GetPoolSnapshotOrSet(ctx, pool)

	rate, err := pool.GetTokenARate(ctx, k.oracleKeeper, &snapshot, tokenInDenom, baseCurrency, k.accountedPoolKeeper)
	if err != nil {
		return sdk.ZeroDec()
	}

	return rate
}

func (k Keeper) GetEdenDenomPrice(ctx sdk.Context, baseCurrency string) math.LegacyDec {
	// Calc ueden / uusdc rate
	edenUsdcRate := k.EstimatePrice(ctx, ptypes.Elys, baseCurrency)
	if edenUsdcRate.IsZero() {
		edenUsdcRate = sdk.OneDec()
	}
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal := int64(6)
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = int64(usdcEntry.Decimals)
		}
		usdcDenomPrice = sdk.NewDecWithPrec(1, usdcDecimal)
	}
	return edenUsdcRate.Mul(usdcDenomPrice)
}
