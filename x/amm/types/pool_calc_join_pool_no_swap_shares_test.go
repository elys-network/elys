// pool_calc_join_pool_no_swap_shares_test.go
package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestMaximalExactRatioJoin(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		pool           *types.Pool
		tokensIn       sdk.Coins
		expectedShares sdkmath.Int
		expectedRem    sdk.Coins
		expectedErrMsg string
	}{
		{
			"successful join with exact ratio",
			&types.Pool{
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(1000)),
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000))},
				},
			},
			sdk.Coins{
				sdk.NewCoin("tokenA", sdkmath.NewInt(100)),
				sdk.NewCoin("tokenB", sdkmath.NewInt(200)),
			},
			sdkmath.NewInt(100),
			sdk.Coins{},
			"",
		},
		{
			"successful join with remaining coins",
			&types.Pool{
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(1000)),
				PoolAssets: []types.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000))},
				},
			},
			sdk.Coins{
				sdk.NewCoin("tokenA", sdkmath.NewInt(150)),
				sdk.NewCoin("tokenB", sdkmath.NewInt(200)),
			},
			sdkmath.NewInt(100),
			sdk.Coins{
				sdk.NewCoin("tokenA", sdkmath.NewInt(50)),
			},
			"",
		},
		{
			"unexpected error due to pool liquidity is zero for denom",
			&types.Pool{
				TotalShares: sdk.NewCoin("shares", sdkmath.NewInt(1000)),
				PoolAssets:  []types.PoolAsset{},
			},
			sdk.Coins{
				sdk.NewCoin("tokenA", sdkmath.NewInt(100)),
			},
			sdkmath.Int{},
			sdk.Coins{},
			"pool liquidity is zero for denom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shares, rem, err := types.MaximalExactRatioJoin(tc.pool, tc.tokensIn)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedShares, shares)
				require.Equal(t, tc.expectedRem, rem)
			}
		})
	}
}
