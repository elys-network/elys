package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestPool_GetAllPoolAssets(t *testing.T) {
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

	result := pool.GetAllPoolAssets()

	require.Equal(t, len(poolAssets), len(result))
	for i := 0; i < len(poolAssets); i++ {
		require.Equal(t, poolAssets[i], result[i])
	}
}
