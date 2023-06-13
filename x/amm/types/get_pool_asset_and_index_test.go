package types_test

import (
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestPool_GetPoolAssetAndIndex(t *testing.T) {
	poolAssets := []*types.PoolAsset{
		{
			Token:  sdk.NewCoin("token1", sdk.NewInt(100)),
			Weight: sdk.NewInt(10),
		},
		{
			Token:  sdk.NewCoin("token2", sdk.NewInt(200)),
			Weight: sdk.NewInt(20),
		},
	}

	pool := types.Pool{
		PoolAssets: poolAssets,
	}

	// Test case 1: Existing PoolAsset
	index, poolAsset, err := pool.GetPoolAssetAndIndex("token1")
	require.NoError(t, err)
	require.Equal(t, 0, index)
	require.Equal(t, poolAssets[0], &poolAsset)

	// Test case 2: Non-existing PoolAsset
	nonExistingDenom := "nonExistingToken"
	_, _, err = pool.GetPoolAssetAndIndex(nonExistingDenom)
	expectedErr := errorsmod.Wrapf(types.ErrDenomNotFoundInPool, fmt.Sprintf(types.FormatNoPoolAssetFoundErrFormat, nonExistingDenom))
	require.EqualError(t, err, expectedErr.Error())

	// Test case 3: Empty denom
	_, _, err = pool.GetPoolAssetAndIndex("")
	require.EqualError(t, err, "you tried to find the PoolAsset with empty denom")
}
