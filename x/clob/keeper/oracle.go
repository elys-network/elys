package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAssetPrice(ctx sdk.Context, denom string) (math.LegacyDec, error) {
	assetInfo, found := k.oracleKeeper.GetAssetInfo(ctx, denom)
	if !found {
		return math.LegacyDec{}, fmt.Errorf("asset info (%s) not found for denom", denom)
	}

	price, found := k.oracleKeeper.GetAssetPrice(ctx, assetInfo.Display)
	if !found {
		return math.LegacyDec{}, fmt.Errorf("asset price not found for denom (%s)", denom)
	}
	return price.Price, nil
}

func (k Keeper) GetDenomPrice(ctx sdk.Context, denom string) (math.LegacyDec, error) {
	price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, denom)
	if price.IsNil() || price.IsZero() {
		return math.LegacyDec{}, fmt.Errorf("denom (%s) price not found", denom)
	}
	return price, nil
}
