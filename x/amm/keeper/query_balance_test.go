package keeper_test

import "github.com/elys-network/elys/v5/x/amm/types"

func (suite *AmmKeeperTestSuite) TestQueryBalance() {
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
				_, err := suite.app.AmmKeeper.Balance(suite.ctx, nil)
				suite.Require().Error(err)
			},
		},
		{
			"invalid address",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.Balance(suite.ctx, &types.QueryBalanceRequest{
					Address: "invalid_address",
				})
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
