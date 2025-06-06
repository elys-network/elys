package keeper_test

import (
	"github.com/elys-network/elys/v6/x/oracle/types"
)

func (suite *KeeperTestSuite) TestGetParams() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	params := types.DefaultParams()

	k.SetParams(ctx, params)
	suite.Require().Equal(params, k.GetParams(ctx))
}
