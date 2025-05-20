package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// GetAssetPriceAndAssetUsdcDenomRatio returns asset price and asset price in denom ratio(price*(10^usdcDecimals-assetDecimals))
// Units are uusdc per base token, example: uusdc per uatom, uusdc per wei, uusdc per satoshi
func (k Keeper) GetAssetPriceAndAssetUsdcDenomRatio(ctx sdk.Context, asset string) (osmomath.BigDec, osmomath.BigDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset info %s not found", asset)
	}
	USDCInfo, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset info %s not found", ptypes.BaseCurrency)
	}

	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", asset)
	}
	if price.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s is zero", asset)
	}
	USDCPrice, found := k.oracleKeeper.GetAssetPrice(ctx, USDCInfo.DisplayName)
	if !found {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", ptypes.BaseCurrency)
	}
	if USDCPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s is zero", ptypes.BaseCurrency)
	}

	if info.Decimal < USDCInfo.Decimals {
		baseUnitMultiplier := utils.Pow10Int64(USDCInfo.Decimals - info.Decimal)
		return price, price.Quo(USDCPrice).MulInt64(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := utils.Pow10Int64(info.Decimal - USDCInfo.Decimals)
		return price, price.Quo(USDCPrice).QuoInt64(baseUnitMultiplier), nil
	}
}

func (k Keeper) ConvertPriceToAssetUsdcDenomRatio(ctx sdk.Context, asset string, price osmomath.BigDec) (osmomath.BigDec, error) {
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
		return price.Quo(USDCPrice).MulInt64(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := utils.Pow10Int64(info.Decimal - USDCInfo.Decimals)
		return price.Quo(USDCPrice).QuoInt64(baseUnitMultiplier), nil
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
