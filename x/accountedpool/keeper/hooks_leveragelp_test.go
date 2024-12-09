// hooks_leveragelp_test.go
package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/app"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestOnLeverageLpPoolEnable(t *testing.T) {
	app := app.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	apk := app.AccountedPoolKeeper

	err := simapp.SetStakingParam(app, ctx)
	require.NoError(t, err)

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	testCases := []struct {
		name           string
		ammPool        ammtypes.Pool
		expectedTokens []sdk.Coin
		expectedErrMsg string
	}{
		{
			"successful pool enable",
			ammtypes.Pool{
				PoolId:      1,
				Address:     addr[0].String(),
				PoolParams:  ammtypes.PoolParams{},
				TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
				PoolAssets: []ammtypes.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000))},
				},
				TotalWeight:       sdkmath.NewInt(100),
				RebalanceTreasury: addr[0].String(),
			},
			[]sdk.Coin{
				sdk.NewCoin("tokenA", sdkmath.NewInt(1000)),
				sdk.NewCoin("tokenB", sdkmath.NewInt(2000)),
			},
			"",
		},
		{
			"pool already exists",
			ammtypes.Pool{
				PoolId:      1,
				Address:     addr[0].String(),
				PoolParams:  ammtypes.PoolParams{},
				TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
				PoolAssets: []ammtypes.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
				},
				TotalWeight:       sdkmath.NewInt(100),
				RebalanceTreasury: addr[0].String(),
			},
			nil,
			"accounted pool already exist",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := apk.OnLeverageLpPoolEnable(ctx, tc.ammPool)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				if tc.expectedTokens != nil {
					accountedPool, found := apk.GetAccountedPool(ctx, tc.ammPool.PoolId)
					require.True(t, found)
					require.Equal(t, tc.expectedTokens, accountedPool.TotalTokens)
				}
			}
		})
	}
}

func TestOnLeverageLpPoolDisable(t *testing.T) {
	app := app.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	apk := app.AccountedPoolKeeper

	err := simapp.SetStakingParam(app, ctx)
	require.NoError(t, err)

	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	testCases := []struct {
		name           string
		setup          func()
		ammPool        ammtypes.Pool
		expectedErrMsg string
	}{
		{
			"accounted pool not found",
			func() {},
			ammtypes.Pool{
				PoolId:      1,
				Address:     addr[0].String(),
				PoolParams:  ammtypes.PoolParams{},
				TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
				PoolAssets: []ammtypes.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
				},
				TotalWeight:       sdkmath.NewInt(100),
				RebalanceTreasury: addr[0].String(),
			},
			types.ErrPoolDoesNotExist.Error(),
		},
		{
			"non-zero non-AMM pool balance",
			func() {
				accountedPool := types.AccountedPool{
					PoolId:           2,
					TotalTokens:      []sdk.Coin{},
					NonAmmPoolTokens: []sdk.Coin{sdk.NewCoin("tokenA", sdkmath.NewInt(500))},
				}

				apk.SetAccountedPool(ctx, accountedPool)
			},
			ammtypes.Pool{
				PoolId:      2,
				Address:     addr[0].String(),
				PoolParams:  ammtypes.PoolParams{},
				TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
				PoolAssets: []ammtypes.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
				},
				TotalWeight:       sdkmath.NewInt(100),
				RebalanceTreasury: addr[0].String(),
			},
			"all positions are not closed; accounted pool have non-zero non amm pool balance left",
		},
		{
			"successful pool disable",
			func() {
				accountedPool := types.AccountedPool{
					PoolId:           1,
					TotalTokens:      []sdk.Coin{},
					NonAmmPoolTokens: []sdk.Coin{},
				}
				accountedPool.TotalTokens = append(accountedPool.TotalTokens, sdk.NewCoin("tokenA", sdkmath.NewInt(1000)))
				accountedPool.TotalTokens = append(accountedPool.TotalTokens, sdk.NewCoin("tokenB", sdkmath.NewInt(2000)))

				apk.SetAccountedPool(ctx, accountedPool)
			},
			ammtypes.Pool{
				PoolId:      1,
				Address:     addr[0].String(),
				PoolParams:  ammtypes.PoolParams{},
				TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
				PoolAssets: []ammtypes.PoolAsset{
					{Token: sdk.NewCoin("tokenA", sdkmath.NewInt(1000))},
					{Token: sdk.NewCoin("tokenB", sdkmath.NewInt(2000))},
				},
				TotalWeight:       sdkmath.NewInt(100),
				RebalanceTreasury: addr[0].String(),
			},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			err := apk.OnLeverageLpPoolDisable(ctx, tc.ammPool)
			if tc.expectedErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErrMsg)
			} else {
				require.NoError(t, err)
				_, found := apk.GetAccountedPool(ctx, tc.ammPool.PoolId)
				require.False(t, found)
			}
		})
	}
}
