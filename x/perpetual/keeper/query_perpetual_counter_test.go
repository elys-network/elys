package keeper_test

import "github.com/elys-network/elys/v7/x/perpetual/types"

func (suite *PerpetualKeeperTestSuite) TestQueryPerpetualCounter_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.PerpetualCounter(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetStatus_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.PerpetualCounter(ctx, &types.PerpetualCounterRequest{Id: 1})

	suite.Require().Nil(err)
}
