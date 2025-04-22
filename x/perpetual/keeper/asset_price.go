package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	return price.Price, nil
}

func (k Keeper) GetAssetPriceInBaseUnits(ctx sdk.Context, asset string) (math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}

	USDCDecimals := uint64(6)
	if info.GetDecimal() < USDCDecimals {
		baseUnitMultiplier := math.LegacyMustNewDecFromStr("10").Power(USDCDecimals - info.GetDecimal())
		return price.Price.Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := math.LegacyMustNewDecFromStr("10").Power(info.GetDecimal() - USDCDecimals)
		return price.Price.Quo(baseUnitMultiplier), nil
	}
}

func (k Keeper) ConvertPriceToBaseUnit(ctx sdk.Context, asset string, price math.LegacyDec) (math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("error converting price to base units, asset info %s not found", asset)
	}

	USDCDecimals := uint64(6)
	if info.GetDecimal() < USDCDecimals {
		baseUnitMultiplier := math.LegacyMustNewDecFromStr("10").Power(USDCDecimals - info.GetDecimal())
		return price.Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := math.LegacyMustNewDecFromStr("10").Power(info.GetDecimal() - USDCDecimals)
		return price.Quo(baseUnitMultiplier), nil
	}
}
