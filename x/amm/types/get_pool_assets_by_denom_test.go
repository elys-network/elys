package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestGetPoolAssetsByDenom(t *testing.T) {
	poolAssets := []types.PoolAsset{
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdk.NewInt(100)},
			Weight: sdk.NewInt(10),
		},
		{
			Token:  sdk.Coin{Denom: "token2", Amount: sdk.NewInt(200)},
			Weight: sdk.NewInt(20),
		},
	}

	// Test case 1: No duplicate pool assets
	poolAssetsByDenom, err := types.GetPoolAssetsByDenom(poolAssets)
	require.NoError(t, err)
	require.Equal(t, 2, len(poolAssetsByDenom))
	require.Equal(t, poolAssets[0], poolAssetsByDenom["token1"])
	require.Equal(t, poolAssets[1], poolAssetsByDenom["token2"])

	// Test case 2: Duplicate pool asset
	duplicatePoolAssets := []types.PoolAsset{
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdk.NewInt(100)},
			Weight: sdk.NewInt(10),
		},
		{
			Token:  sdk.Coin{Denom: "token1", Amount: sdk.NewInt(200)},
			Weight: sdk.NewInt(20),
		},
	}
	_, err = types.GetPoolAssetsByDenom(duplicatePoolAssets)
	expectedErr := fmt.Errorf(types.FormatRepeatingPoolAssetsNotAllowedErrFormat, "token1")
	require.EqualError(t, err, expectedErr.Error())
}
