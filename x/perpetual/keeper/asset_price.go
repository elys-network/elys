package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAssetPriceByDenom(ctx sdk.Context, asset string) (sdk.Dec, error) {
	info, found := k.oracleKeeper.GetAssetInfo(ctx, asset)
	if !found {
		return sdk.ZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return sdk.ZeroDec(), fmt.Errorf("asset price %s not found", asset)
	}
	return price.Price, nil
}
