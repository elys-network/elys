package keeper_test

import (
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCheckMaxOpenPositions() {
	params := types.DefaultParams()
	params.MaxOpenPositions = 2
	err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"Open Positions Below Max",
			"",
			func() {
				suite.app.PerpetualKeeper.SetOpenMTPCount(suite.ctx, 1)
			},
		},
		{
			"Open Positions Equal Max",
			types.ErrMaxOpenPositions.Error(),
			func() {
				suite.app.PerpetualKeeper.SetOpenMTPCount(suite.ctx, 2)
			},
		},
		{
			"Open Positions Exceed Max",
			types.ErrMaxOpenPositions.Error(),
			func() {
				suite.app.PerpetualKeeper.SetOpenMTPCount(suite.ctx, 3)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err = suite.app.PerpetualKeeper.CheckMaxOpenPositions(suite.ctx)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
