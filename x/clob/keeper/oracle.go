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
		return math.LegacyDec{}, fmt.Errorf("asset info (%s) not found for denom", denom)
	}

	price, found := k.oracleKeeper.GetAssetPrice(ctx, assetInfo.Display)
	if !found {
		return math.LegacyDec{}, fmt.Errorf("asset price not found for denom (%s)", denom)
	}
	if price.LTE(math.LegacyZeroDec()) || price.IsNil() {
		return math.LegacyDec{}, fmt.Errorf("asset price (%s) is invalid", price.String())
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
