package keeper_test

import (
	"github.com/elys-network/elys/v6/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestParamsQuery() {
	params := types.DefaultParams()
	suite.app.MasterchefKeeper.SetParams(suite.ctx, params)

	response, err := suite.app.MasterchefKeeper.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
}
