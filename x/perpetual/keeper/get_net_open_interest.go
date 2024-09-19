package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetNetOpenInterest(pool types.Pool) math.Int {
	assetLiabilitiesLong := sdk.ZeroInt()
	assetLiabilitiesShort := sdk.ZeroInt()

	for _, asset := range pool.PoolAssetsLong {
		assetLiabilitiesLong = assetLiabilitiesLong.Add(asset.Liabilities)
	}

	for _, asset := range pool.PoolAssetsShort {
		assetLiabilitiesShort = assetLiabilitiesShort.Add(asset.Liabilities)
	}

	netOpenInterest := assetLiabilitiesLong.Sub(assetLiabilitiesShort)
	return netOpenInterest
}
