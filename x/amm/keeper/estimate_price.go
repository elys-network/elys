package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
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
	usdcDenomPrice, usdcDecimal := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal = 6
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = usdcEntry.Decimals
		}
		usdcDenomPrice = math.LegacyNewDec(1)
	}
	return edenUsdcRate.Mul(usdcDenomPrice).Quo(oraclekeeper.Pow10(usdcDecimal))
}

func (k Keeper) GetTokenPrice(ctx sdk.Context, tokenInDenom, baseCurrency string) math.LegacyDec {
	oraclePrice, oracleDecimal := k.oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if !oraclePrice.IsZero() {
		return oraclePrice.Quo(oraclekeeper.Pow10(oracleDecimal))
	}

	// Calc tokenIn / uusdc rate
	tokenUsdcRate := k.EstimatePrice(ctx, tokenInDenom, baseCurrency)
	usdcDenomPrice, usdcDecimal := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	if usdcDenomPrice.IsZero() {
		usdcDecimal = 6
		usdcEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if found {
			usdcDecimal = usdcEntry.Decimals
		}
		// usdcDenomPrice = math.LegacyNewDecWithPrec(1, usdcDecimal)
		usdcDenomPrice = math.LegacyNewDec(1)
	}
	return tokenUsdcRate.Mul(usdcDenomPrice).Quo(oraclekeeper.Pow10(usdcDecimal))
}
