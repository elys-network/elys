package keeper

import (
	"cosmossdk.io/math"
	"fmt"
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
