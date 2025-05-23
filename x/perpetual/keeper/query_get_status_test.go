package keeper_test

import "github.com/elys-network/elys/v5/x/perpetual/types"

func (suite *PerpetualKeeperTestSuite) TestQueryGetStatus_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetStatus(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetStatus_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetStatus(ctx, &types.StatusRequest{})

	suite.Require().Nil(err)
}
