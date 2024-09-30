package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetNetOpenInterest(pool types.Pool) sdkmath.Int {
	assetLiabilitiesLong := sdkmath.ZeroInt()
	assetLiabilitiesShort := sdkmath.ZeroInt()

	for _, asset := range pool.PoolAssetsLong {
		assetLiabilitiesLong = assetLiabilitiesLong.Add(asset.Liabilities)
	}

	for _, asset := range pool.PoolAssetsShort {
		assetLiabilitiesShort = assetLiabilitiesShort.Add(asset.Liabilities)
	}

	netOpenInterest := assetLiabilitiesLong.Sub(assetLiabilitiesShort)
	return netOpenInterest
}
