package keeper_test

import "github.com/elys-network/elys/x/perpetual/types"

func (suite *PerpetualKeeperTestSuite) TestQueryGetAllToPay_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetAllToPay(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetAllToPay_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetAllToPay(ctx, &types.QueryGetAllToPayRequest{})

	suite.Require().Nil(err)
}
