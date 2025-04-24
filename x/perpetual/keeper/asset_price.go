package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracleKeeper "github.com/elys-network/elys/x/oracle/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// returns asset price and and asset price in denom ratio(price*(10^usdcDecimals-assetDecimals))
func (k Keeper) GetAssetPriceAndAssetUsdcDenomRatio(ctx sdk.Context, asset string) (math.LegacyDec, math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), fmt.Errorf("asset info %s not found", asset)
	}
	USDCInfo, found := k.oracleKeeper.GetAssetInfo(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), fmt.Errorf("asset info %s not found", ptypes.BaseCurrency)
	}

	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	USDCPrice, found := k.oracleKeeper.GetAssetPrice(ctx, USDCInfo.Display)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", ptypes.BaseCurrency)
	}

	if info.Decimal < USDCInfo.Decimal {
		baseUnitMultiplier := oracleKeeper.Pow10(USDCInfo.Decimal - info.Decimal)
		return price.Price, price.Price.Quo(USDCPrice.Price).Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := oracleKeeper.Pow10(info.Decimal - USDCInfo.Decimal)
		return price.Price, price.Price.Quo(USDCPrice.Price).Quo(baseUnitMultiplier), nil
	}
}

func (k Keeper) ConvertPriceToAssetUsdcDenomRatio(ctx sdk.Context, asset string, price math.LegacyDec) (math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("error converting price to base units, asset info %s not found", asset)
	}
	USDCInfo, found := k.oracleKeeper.GetAssetInfo(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("error converting price to base units, asset info %s not found", ptypes.BaseCurrency)
	}
	USDCPrice, found := k.oracleKeeper.GetAssetPrice(ctx, USDCInfo.Display)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("error converting price to base units, asset price %s not found", ptypes.BaseCurrency)
	}
	if info.Decimal < USDCInfo.Decimal {
		baseUnitMultiplier := oracleKeeper.Pow10(USDCInfo.Decimal - info.Decimal)
		return price.Quo(USDCPrice.Price).Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := oracleKeeper.Pow10(info.Decimal - USDCInfo.Decimal)
		return price.Quo(USDCPrice.Price).Quo(baseUnitMultiplier), nil
	}
}
