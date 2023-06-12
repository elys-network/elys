package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestGetPoolAssetByDenom(t *testing.T) {
	poolAssets := []*types.PoolAsset{
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdk.NewInt(100)},
			Weight: sdk.NewInt(10),
		},
		{
			Token:  sdk.Coin{Denom: "token2", Amount: sdk.NewInt(200)},
			Weight: sdk.NewInt(20),
		},
	}

	// Test case 1: Existing PoolAsset
	asset, found := types.GetPoolAssetByDenom(poolAssets, "token1")
	require.True(t, found)
	require.Equal(t, poolAssets[0], asset)

	// Test case 2: Non-existing PoolAsset
	_, found = types.GetPoolAssetByDenom(poolAssets, "nonExistingToken")
	require.False(t, found)
}
