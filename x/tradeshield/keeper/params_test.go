package keeper_test

import (
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestGetParams() {
	params := types.DefaultParams()
	suite.app.TradeshieldKeeper.SetParams(suite.ctx, &params)

	suite.Require().EqualValues(params, suite.app.TradeshieldKeeper.GetParams(suite.ctx))
}
