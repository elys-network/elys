package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"

	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

// GetPerpetualPoolBalances
func TestKeeper_GetPerpetualPoolBalances(t *testing.T) {
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: math.NewInt(100),
				},
			},
		},
	}
	perpetualPool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("11.5"))
	perpetualPool.PoolAssetsLong = []types.PoolAsset{
		{
			AssetDenom:            "testAsset",
			Liabilities:           math.NewInt(100),
			Custody:               math.NewInt(100),
			TakeProfitLiabilities: math.NewInt(100),
			TakeProfitCustody:     math.NewInt(100),
		},
	}
	perpetualPool.PoolAssetsShort = []types.PoolAsset{
		{
			AssetDenom:            "testAsset",
			Liabilities:           math.NewInt(100),
			Custody:               math.NewInt(100),
			TakeProfitLiabilities: math.NewInt(100),
			TakeProfitCustody:     math.NewInt(100),
		},
	}

	liabilities, custody, _, _ := perpetualPool.GetPerpetualPoolBalances("testAsset")

	require.Equal(t, math.NewInt(200), liabilities)
	require.Equal(t, math.NewInt(200), custody)
}
