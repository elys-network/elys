package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestSlippageTrack() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"request with nil",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.SlippageTrack(suite.ctx, nil)
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
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				_, err := suite.app.AmmKeeper.SlippageTrack(suite.ctx, &types.QuerySlippageTrackRequest{
					PoolId: pool.PoolId,
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

func (suite *AmmKeeperTestSuite) TestSlippageTrackAll() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"request with nil",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.SlippageTrackAll(suite.ctx, nil)
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

				_, err := suite.app.AmmKeeper.SlippageTrackAll(suite.ctx, &types.QuerySlippageTrackAllRequest{})
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
