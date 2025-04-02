package keeper_test

import (
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestGetPerpetualMarket() {

	market := types.PerpetualMarket{
		Id: 1,
	}
	suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
	// Define test cases
	testCases := []struct {
		name           string
		id             uint64
		expectedErrMsg string
	}{
		{
			"market not found",
			2,
			types.ErrPerpetualMarketNotFound.Error(),
		},
		{
			"market found",
			1,
			"",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			res, err := suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, tc.id)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.id, res.Id)
			}
		})
	}
}
