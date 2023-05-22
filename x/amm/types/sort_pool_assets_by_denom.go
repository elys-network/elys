package types

import (
	"sort"
	"strings"
)

// sortPoolAssetsByDenom sorts pool assets in place, by weight.
func sortPoolAssetsByDenom(assets []*PoolAsset) {
	sort.Slice(assets, func(i, j int) bool {
		PoolAssetA := assets[i]
		PoolAssetB := assets[j]

		return strings.Compare(PoolAssetA.Token.Denom, PoolAssetB.Token.Denom) == -1
	})
}
