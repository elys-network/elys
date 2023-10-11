package keeper

import (
	"github.com/elys-network/elys/x/leveragelp/types"
)

// Get Assets Index
func (k Keeper) GetMTPAssetIndex(mtp *types.MTP, collateralAsset string) int {
	collateralIndex := -1
	for i, asset := range mtp.CollateralAssets {
		if asset == collateralAsset {
			collateralIndex = i
			break
		}
	}

	return collateralIndex
}
