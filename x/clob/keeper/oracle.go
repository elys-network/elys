package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetAssetPriceFromDenom(ctx sdk.Context, denom string) (math.LegacyDec, error) {
	assetInfo, found := k.oracleKeeper.GetAssetInfo(ctx, denom)
	if !found {
		return math.LegacyDec{}, WrapPriceNotFoundError(denom, "asset info lookup")
	}

	price, found := k.oracleKeeper.GetAssetPrice(ctx, assetInfo.Display)
	if !found {
		return math.LegacyDec{}, WrapPriceNotFoundError(denom, "oracle")
	}
	if price.LTE(math.LegacyZeroDec()) || price.IsNil() {
		return math.LegacyDec{}, WrapInvalidPriceError(price, "price must be positive")
	}
	return price, nil
}

func (k Keeper) GetDenomPrice(ctx sdk.Context, denom string) (osmomath.BigDec, math.LegacyDec, error) {
	price := k.oracleKeeper.GetDenomPrice(ctx, denom)
	if price.IsNil() || price.IsZero() {
		return osmomath.BigDec{}, math.LegacyDec{}, fmt.Errorf("denom (%s) price not found", denom)
	}
	// No major benefit to use 36 decimal places in clob as everything is synthetic except margin amount
	return price, price.Dec(), nil
}
