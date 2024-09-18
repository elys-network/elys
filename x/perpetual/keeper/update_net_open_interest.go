package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateNetOpenInterest(ctx sdk.Context, pool *types.Pool) error {
	assetLiabilitiesLong := sdk.ZeroInt()
	assetLiabilitiesShort := sdk.ZeroInt()

	for _, asset := range pool.PoolAssetsLong {
		assetLiabilitiesLong = assetLiabilitiesLong.Add(asset.Liabilities)
	}

	for _, asset := range pool.PoolAssetsShort {
		assetLiabilitiesShort = assetLiabilitiesShort.Add(asset.Liabilities)
	}

	pool.NetOpenInterest = assetLiabilitiesLong.Sub(assetLiabilitiesShort)
	return nil
}
