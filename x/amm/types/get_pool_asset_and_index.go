package types

import (
	"fmt"
	"sort"
	"strings"

	errorsmod "cosmossdk.io/errors"
)

const (
	formatNoPoolAssetFoundErrFormat = "can't find the PoolAsset (%s)"
)

// Returns a pool asset, and its index. If err != nil, then the index will be valid.
func (p Pool) getPoolAssetAndIndex(denom string) (int, PoolAsset, error) {
	if denom == "" {
		return -1, PoolAsset{}, fmt.Errorf("you tried to find the PoolAsset with empty denom")
	}

	if len(p.PoolAssets) == 0 {
		return -1, PoolAsset{}, errorsmod.Wrapf(ErrDenomNotFoundInPool, fmt.Sprintf(formatNoPoolAssetFoundErrFormat, denom))
	}

	i := sort.Search(len(p.PoolAssets), func(i int) bool {
		PoolAssetA := p.PoolAssets[i]

		compare := strings.Compare(PoolAssetA.Token.Denom, denom)
		return compare >= 0
	})

	if i < 0 || i >= len(p.PoolAssets) {
		return -1, PoolAsset{}, errorsmod.Wrapf(ErrDenomNotFoundInPool, fmt.Sprintf(formatNoPoolAssetFoundErrFormat, denom))
	}

	if p.PoolAssets[i].Token.Denom != denom {
		return -1, PoolAsset{}, errorsmod.Wrapf(ErrDenomNotFoundInPool, fmt.Sprintf(formatNoPoolAssetFoundErrFormat, denom))
	}

	return i, *p.PoolAssets[i], nil
}
