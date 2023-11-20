package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/stretchr/testify/require"
)

// GetMarginPoolBalances
func TestKeeper_GetMarginPoolBalances(t *testing.T) {
	marginPool := types.NewPool(1)
	marginPool.PoolAssetsLong = append(marginPool.PoolAssetsLong, types.PoolAsset{
		AssetDenom:   "testAsset",
		AssetBalance: sdk.NewInt(100),
		Liabilities:  sdk.NewInt(100),
		Custody:      sdk.NewInt(100),
	})
	marginPool.PoolAssetsShort = append(marginPool.PoolAssetsShort, types.PoolAsset{
		AssetDenom:   "testAsset",
		AssetBalance: sdk.NewInt(100),
		Liabilities:  sdk.NewInt(100),
		Custody:      sdk.NewInt(100),
	})

	assetBalance, liabilities, custody := types.GetMarginPoolBalances(marginPool, "testAsset")

	require.Equal(t, sdk.NewInt(200), assetBalance)
	require.Equal(t, sdk.NewInt(200), liabilities)
	require.Equal(t, sdk.NewInt(200), custody)
}
