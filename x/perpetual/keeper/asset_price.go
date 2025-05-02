package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (osmomath.BigDec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", asset)
	}
	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("asset price %s not found", asset)
	}
	return price, nil
}
