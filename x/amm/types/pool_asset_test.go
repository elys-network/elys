package types_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
)

func (suite *TestSuite) TestPoolAssetValidate() {
	suite.SetupTest()
	poolAsset := types.PoolAsset{
		Token:                  sdk.Coin{},
		Weight:                 sdkmath.Int{},
		ExternalLiquidityRatio: sdkmath.LegacyDec{},
	}
	for _, tc := range []struct {
		name     string
		errorMsg string
		function func()
	}{
		{
			name:     "token is nil",
			errorMsg: "invalid pool asset token",
			function: func() {},
		},
		{
			name:     "invalid token denom",
			errorMsg: "invalid pool asset token",
			function: func() {
				poolAsset.Token = sdk.Coin{"%%%", sdkmath.OneInt()}
			},
		},
		{
			name:     "token is -ve",
			errorMsg: "invalid pool asset token",
			function: func() {
				poolAsset.Token = sdk.Coin{"uatom", sdkmath.OneInt().MulRaw(-1)}
			},
		},
		{
			name:     "weight is nil",
			errorMsg: "invalid pool asset weight",
			function: func() {
				poolAsset.Token = sdk.Coin{"uatom", sdkmath.OneInt()}
			},
		},
		{
			name:     "weight is -ve",
			errorMsg: "invalid pool asset weight",
			function: func() {
				poolAsset.Weight = sdkmath.NewInt(-1)
			},
		},
		{
			name:     "external liquidity ratio is nil",
			errorMsg: "invalid external liquidity ratio",
			function: func() {
				poolAsset.Weight = sdkmath.OneInt()
			},
		},
		{
			name:     "external liquidity ratio is < 1",
			errorMsg: "invalid external liquidity ratio",
			function: func() {
				poolAsset.ExternalLiquidityRatio = sdkmath.LegacyNewDec(-1)
			},
		},
		{
			name:     "success",
			errorMsg: "",
			function: func() {
				poolAsset.ExternalLiquidityRatio = sdkmath.LegacyNewDec(1)
			},
		},
	} {
		suite.Run(tc.name, func() {
			tc.function()
			err := poolAsset.Validate()
			if err != nil {
				suite.Require().Contains(err.Error(), tc.errorMsg)
			} else {
				suite.Require().Equal(tc.errorMsg, "")
			}
		})
	}
}
