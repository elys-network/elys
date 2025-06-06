package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/utils"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// GetAssetPriceAndAssetUsdcDenomRatio returns asset price and asset price in denom ratio(price*(10^usdcDecimals-assetDecimals))
// Units are uusdc per base token, example: uusdc per uatom, uusdc per wei, uusdc per satoshi
func (k Keeper) GetAssetPriceAndAssetUsdcDenomRatio(ctx sdk.Context, asset string) (math.LegacyDec, osmomath.BigDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset info %s not found", asset)
	}
	USDCInfo, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset info %s not found", ptypes.BaseCurrency)
	}

	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return math.LegacyZeroDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", asset)
	}
	if price.IsZero() {
		return math.LegacyZeroDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s is zero", asset)
	}
	USDCPrice, found := k.oracleKeeper.GetAssetPrice(ctx, USDCInfo.DisplayName)
	if !found {
		return math.LegacyZeroDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", ptypes.BaseCurrency)
	}
	if USDCPrice.IsZero() {
		return math.LegacyZeroDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s is zero", ptypes.BaseCurrency)
	}

	if info.Decimal < USDCInfo.Decimals {
		baseUnitMultiplier := utils.Pow10Int64(USDCInfo.Decimals - info.Decimal)
		return price.Dec(), price.Quo(USDCPrice).MulInt64(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := utils.Pow10Int64(info.Decimal - USDCInfo.Decimals)
		return price.Dec(), price.Quo(USDCPrice).QuoInt64(baseUnitMultiplier), nil
	}
}

// ConvertPriceToAssetUsdcDenomRatio converts usd per atom, usd per eth, usd per btc to uusdc per base token, example: uusdc per uatom, uusdc per wei, uusdc per satoshi
func (k Keeper) ConvertPriceToAssetUsdcDenomRatio(ctx sdk.Context, asset string, price math.LegacyDec) (osmomath.BigDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("error converting price to base units, asset info %s not found", asset)
	}
	USDCInfo, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("error converting price to base units, asset info %s not found", ptypes.BaseCurrency)
	}
	USDCPrice, found := k.oracleKeeper.GetAssetPrice(ctx, USDCInfo.DisplayName)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("error converting price to base units, asset price %s not found", ptypes.BaseCurrency)
	}
	if info.Decimal < USDCInfo.Decimals {
		baseUnitMultiplier := utils.Pow10Int64(USDCInfo.Decimals - info.Decimal)
		return osmomath.BigDecFromDec(price).Quo(USDCPrice).MulInt64(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := utils.Pow10Int64(info.Decimal - USDCInfo.Decimals)
		return osmomath.BigDecFromDec(price).Quo(USDCPrice).QuoInt64(baseUnitMultiplier), nil
	}
}

func (k Keeper) GetDenomPrice(ctx sdk.Context, denom string) (osmomath.BigDec, error) {
	price := k.oracleKeeper.GetDenomPrice(ctx, denom)
	if price.IsZero() {
		return osmomath.ZeroBigDec(), fmt.Errorf("denom price %s is zero", denom)
	}
	return price, nil
}

func (k Keeper) GetDenomDecimal(ctx sdk.Context, denom string) (uint64, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, denom)
	if !found {
		return 0, fmt.Errorf("asset info %s not found", denom)
	}
	return info.Decimal, nil
}

// ConvertDenomRatioPriceToUSDPrice price -  units are uusdc per uatom, uusdc per wei, usd per sat
// gets converted to usd per atom, usd per eth and usd per btc
func (k Keeper) ConvertDenomRatioPriceToUSDPrice(ctx sdk.Context, denomRatioPrice osmomath.BigDec, tradingAssetDenom string) (math.LegacyDec, error) {
	// units are uusdc per uatom, uusdc per wei, usd per sat
	USDCInfo, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("asset info %s not found", ptypes.BaseCurrency)
	}
	baseCurrencyDenomPrice, err := k.GetDenomPrice(ctx, USDCInfo.Denom)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	// Now the units are usd per uatom, usd per wei, usd per sat
	denomRatioPrice = denomRatioPrice.Mul(baseCurrencyDenomPrice)

	decimal, err := k.GetDenomDecimal(ctx, tradingAssetDenom)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	// Multiply by 10^decimal of taring asset
	denomRatioPrice = denomRatioPrice.MulInt64(utils.Pow10Int64(decimal))
	return denomRatioPrice.Dec(), nil
}
