package keeper_test

import (
	"github.com/elys-network/elys/v5/x/estaking/types"
)

func (suite *EstakingKeeperTestSuite) TestParams() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"get params",
			func() {
				suite.ResetSuite()

				params := types.DefaultParams()

				suite.app.EstakingKeeper.SetParams(suite.ctx, params)
			},
			func() {
				params := types.DefaultParams()

				suite.Require().EqualValues(params, suite.app.EstakingKeeper.GetParams(suite.ctx))
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
