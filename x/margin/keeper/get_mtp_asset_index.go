package keeper

import (
	"github.com/elys-network/elys/x/margin/types"
)

// Get Assets Index
func (k Keeper) GetMTPAssetIndex(mtp *types.MTP, collateralAsset string, borrowAsset string) (int, int) {
	collateralIndex := -1
	borrowIndex := -1
	for i, asset := range mtp.Collaterals {
		if asset.Denom == collateralAsset {
			collateralIndex = i
			break
		}
	}

	for i, asset := range mtp.Custodies {
		if asset.Denom == borrowAsset {
			borrowIndex = i
			break
		}
	}

	return collateralIndex, borrowIndex
}
