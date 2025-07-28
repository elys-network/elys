package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/v7/x/amm/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestInRouteByDenom() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"nil request",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.InRouteByDenom(suite.ctx, nil)
				suite.Require().Error(err)
			},
		},
		{
			"base currency does not exist in asset profile",
			func() {
				suite.ResetSuite()
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.BaseCurrency)
			},
			func() {
				_, err := suite.app.AmmKeeper.InRouteByDenom(suite.ctx, &types.QueryInRouteByDenomRequest{})
				suite.Require().Error(err)
			},
		},
		{
			"pool does not exist",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.InRouteByDenom(suite.ctx, &types.QueryInRouteByDenomRequest{
					DenomIn:  ptypes.ATOM,
					DenomOut: ptypes.ATOM,
				})
				suite.Require().Error(err)
			},
		},
		{
			"valid request",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := math.NewInt(100000000000)
				_ = suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				_, err := suite.app.AmmKeeper.InRouteByDenom(suite.ctx, &types.QueryInRouteByDenomRequest{
					DenomIn:  ptypes.BaseCurrency,
					DenomOut: ptypes.ATOM,
				})
				suite.Require().NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
