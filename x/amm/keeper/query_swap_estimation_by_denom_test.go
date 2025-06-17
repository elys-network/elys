package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestQuerySwapEstimationByDenom() {
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
				_, err := suite.app.AmmKeeper.SwapEstimationByDenom(suite.ctx, nil)
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
				_, err := suite.app.AmmKeeper.SwapEstimationByDenom(suite.ctx, &types.QuerySwapEstimationByDenomRequest{})
				suite.Require().Error(err)
			},
		},
		{
			"atom does not exist in asset profile",
			func() {
				suite.ResetSuite()
				suite.app.AssetprofileKeeper.RemoveEntry(suite.ctx, ptypes.ATOM)
			},
			func() {
				_, err := suite.app.AmmKeeper.SwapEstimationByDenom(suite.ctx, &types.QuerySwapEstimationByDenomRequest{
					DenomIn: ptypes.ATOM,
				})
				suite.Require().Error(err)
			},
		},
		{
			"pool does not exist",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				_, err := suite.app.AmmKeeper.SwapEstimationByDenom(suite.ctx, &types.QuerySwapEstimationByDenomRequest{
					Amount:   sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					DenomIn:  ptypes.BaseCurrency,
					DenomOut: ptypes.ATOM,
					Address:  "",
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

				_, err := suite.app.AmmKeeper.SwapEstimationByDenom(suite.ctx, &types.QuerySwapEstimationByDenomRequest{
					Amount:   sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					DenomIn:  ptypes.BaseCurrency,
					DenomOut: ptypes.ATOM,
					Address:  "",
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
