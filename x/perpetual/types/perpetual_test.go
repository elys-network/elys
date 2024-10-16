package types_test

import (
	"cosmossdk.io/math"
	"testing"

	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

// GetPerpetualPoolBalances
func TestKeeper_GetPerpetualPoolBalances(t *testing.T) {
	perpetualPool := types.NewPool(1)
	perpetualPool.PoolAssetsLong = append(perpetualPool.PoolAssetsLong, types.PoolAsset{
		AssetDenom:            "testAsset",
		Liabilities:           math.NewInt(100),
		Custody:               math.NewInt(100),
		TakeProfitLiabilities: math.NewInt(100),
		TakeProfitCustody:     math.NewInt(100),
	})
	perpetualPool.PoolAssetsShort = append(perpetualPool.PoolAssetsShort, types.PoolAsset{
		AssetDenom:            "testAsset",
		Liabilities:           math.NewInt(100),
		Custody:               math.NewInt(100),
		TakeProfitLiabilities: math.NewInt(100),
		TakeProfitCustody:     math.NewInt(100),
	})

	liabilities, custody, _, _ := perpetualPool.GetPerpetualPoolBalances("testAsset")

	require.Equal(t, math.NewInt(200), liabilities)
	require.Equal(t, math.NewInt(200), custody)
}
