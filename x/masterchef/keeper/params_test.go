package keeper_test

import (
	"github.com/elys-network/elys/v4/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestGetParams() {
	params := types.DefaultParams()

	suite.app.MasterchefKeeper.SetParams(suite.ctx, params)

	suite.Require().EqualValues(params, suite.app.MasterchefKeeper.GetParams(suite.ctx))
}
