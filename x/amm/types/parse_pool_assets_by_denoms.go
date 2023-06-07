package types

import (
	"fmt"
)

func (p Pool) parsePoolAssetsByDenoms(tokenADenom, tokenBDenom string) (
	Aasset *PoolAsset, Basset *PoolAsset, err error,
) {
	Aasset, found1 := getPoolAssetByDenom(p.PoolAssets, tokenADenom)
	Basset, found2 := getPoolAssetByDenom(p.PoolAssets, tokenBDenom)

	if !found1 {
		return &PoolAsset{}, &PoolAsset{}, fmt.Errorf("(%s) does not exist in the pool", tokenADenom)
	}
	if !found2 {
		return &PoolAsset{}, &PoolAsset{}, fmt.Errorf("(%s) does not exist in the pool", tokenBDenom)
	}
	return Aasset, Basset, nil
}
