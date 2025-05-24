package keeper_test

import (
	"github.com/elys-network/elys/v5/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestParamsQuery() {

	params := types.DefaultParams()
	suite.app.TradeshieldKeeper.SetParams(suite.ctx, &params)

	response, err := suite.app.TradeshieldKeeper.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
}
