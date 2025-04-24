package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracleKeeper "github.com/elys-network/elys/x/oracle/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// returns asset price and and asset price in denom ratio(price*(10^usdcDecimals-assetDecimals))
func (k Keeper) GetAssetPriceAndAssetUsdcDenomRatio(ctx sdk.Context, asset string) (math.LegacyDec, math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	USDCEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	if info.Decimal < USDCEntry.Decimals {
		baseUnitMultiplier := oracleKeeper.Pow10(USDCEntry.Decimals - info.Decimal)
		return price.Price, price.Price.Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := oracleKeeper.Pow10(info.Decimal - USDCEntry.Decimals)
		return price.Price, price.Price.Quo(baseUnitMultiplier), nil
	}
}

func (k Keeper) ConvertPriceToAssetUsdcDenomRatio(ctx sdk.Context, asset string, price math.LegacyDec) (math.LegacyDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return math.LegacyZeroDec(), fmt.Errorf("error converting price to base units, asset info %s not found", asset)
	}
	USDCEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return math.LegacyZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	if info.Decimal < USDCEntry.Decimals {
		baseUnitMultiplier := oracleKeeper.Pow10(USDCEntry.Decimals - info.Decimal)
		return price.Mul(baseUnitMultiplier), nil
	} else {
		baseUnitMultiplier := oracleKeeper.Pow10(info.Decimal - USDCEntry.Decimals)
		return price.Quo(baseUnitMultiplier), nil
	}
}
