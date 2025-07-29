package keeper_test

import (
	"github.com/elys-network/elys/v7/x/estaking/types"
)

func (suite *EstakingKeeperTestSuite) TestParamsQuery() {
	params := types.DefaultParams()
	suite.app.EstakingKeeper.SetParams(suite.ctx, params)

	response, err := suite.app.EstakingKeeper.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryParamsResponse{Params: params}, response)
}
