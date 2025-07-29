package keeper_test

import (
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (suite *AmmKeeperTestSuite) TestQuery() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"params",
			func() {
				suite.ResetSuite()

				params := types.DefaultParams()
				suite.app.AmmKeeper.SetParams(suite.ctx, params)
			},
			func() {
				params := types.DefaultParams()
				params.BaseAssets = nil
				params.AllowedUpfrontSwapMakers = nil
				response, err := suite.app.AmmKeeper.Params(suite.ctx, &types.QueryParamsRequest{})
				suite.Require().NoError(err)
				suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
			},
		},
		{
			"params with nil request",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.Params(suite.ctx, nil)
				suite.Require().Error(err)
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
