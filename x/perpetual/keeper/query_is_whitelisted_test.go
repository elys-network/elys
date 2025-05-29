package keeper_test

import (
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestQueryIsWhitelisted_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.IsWhitelisted(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryIsWhitelisted_ErrAccAddressFromBech32() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.IsWhitelisted(ctx, &types.IsWhitelistedRequest{
		Address: "error",
	})

	suite.Require().ErrorContains(err, "invalid bech32 string length 5")
}

func (suite *PerpetualKeeperTestSuite) TestQueryIsWhitelisted_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	response, err := k.IsWhitelisted(ctx, &types.IsWhitelistedRequest{
		Address: sample.AccAddress(),
	})

	suite.Require().Equal(false, response.IsWhitelisted)
	suite.Require().Nil(err)
}
