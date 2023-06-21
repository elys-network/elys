package types

import (
	"fmt"
)

// GetPoolAssetsByDenom return a mapping from pool asset
// denom to the pool asset itself. There must be no duplicates.
// Returns error, if any found.
func GetPoolAssetsByDenom(poolAssets []PoolAsset) (map[string]PoolAsset, error) {
	poolAssetsByDenom := make(map[string]PoolAsset)
	for _, poolAsset := range poolAssets {
		_, ok := poolAssetsByDenom[poolAsset.Token.Denom]
		if ok {
			return nil, fmt.Errorf(FormatRepeatingPoolAssetsNotAllowedErrFormat, poolAsset.Token.Denom)
		}

		poolAssetsByDenom[poolAsset.Token.Denom] = poolAsset
	}
	return poolAssetsByDenom, nil
}
