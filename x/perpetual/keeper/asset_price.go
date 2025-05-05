package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// returns asset price and and asset price in denom ratio(price*(10^usdcDecimals-assetDecimals))
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
	USDCPrice, found := k.oracleKeeper.GetAssetPrice(ctx, USDCInfo.DisplayName)
	if !found {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", ptypes.BaseCurrency)
	}

	if info.Decimal < USDCInfo.Decimals {
		baseUnitMultiplier := utils.Pow10(USDCInfo.Decimals - info.Decimal)
		return price, price.Quo(USDCPrice).Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := utils.Pow10(info.Decimal - USDCInfo.Decimals)
		return price, price.Quo(USDCPrice).Quo(baseUnitMultiplier), nil
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
		baseUnitMultiplier := utils.Pow10(USDCInfo.Decimals - info.Decimal)
		return price.Quo(USDCPrice).Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := utils.Pow10(info.Decimal - USDCInfo.Decimals)
		return price.Quo(USDCPrice).Quo(baseUnitMultiplier), nil
	}
}
