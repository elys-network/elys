package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

// GetPerpetualPoolBalances
func TestKeeper_GetPerpetualPoolBalances(t *testing.T) {
	perpetualPool := types.NewPool(1)
	perpetualPool.PoolAssetsLong = append(perpetualPool.PoolAssetsLong, types.PoolAsset{
		AssetDenom:   "testAsset",
		AssetBalance: sdk.NewInt(100),
		Liabilities:  sdk.NewInt(100),
		Custody:      sdk.NewInt(100),
	})
	perpetualPool.PoolAssetsShort = append(perpetualPool.PoolAssetsShort, types.PoolAsset{
		AssetDenom:   "testAsset",
		AssetBalance: sdk.NewInt(100),
		Liabilities:  sdk.NewInt(100),
		Custody:      sdk.NewInt(100),
	})

	assetBalance, liabilities, custody := types.GetPerpetualPoolBalances(perpetualPool, "testAsset")

	require.Equal(t, sdk.NewInt(200), assetBalance)
	require.Equal(t, sdk.NewInt(200), liabilities)
	require.Equal(t, sdk.NewInt(200), custody)
}
