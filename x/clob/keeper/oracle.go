package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (math.LegacyDec, error) {
	price, found := k.oracleKeeper.GetAssetPrice(ctx, asset)
	if !found {
		return math.LegacyDec{}, fmt.Errorf("asset (%s) price not found", asset)
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
