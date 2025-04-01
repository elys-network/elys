package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) math.Dec {
	price, found := k.oracleKeeper.GetAssetPrice(ctx, asset)
	if !found {
		return math.Dec{}
	}
	result, err := math.DecFromLegacyDec(price.Price)
	if err != nil {
		panic(err)
	}
	return result
}

func (k Keeper) GetDenomPrice(ctx sdk.Context, denom string) (math.Dec, error) {
	return math.DecFromLegacyDec(k.oracleKeeper.GetAssetPriceFromDenom(ctx, denom))
}
