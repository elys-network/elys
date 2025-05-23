package keeper_test

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestQueryGetWhiteList_InvalidRequest() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetWhitelist(ctx, nil)

	suite.Require().ErrorContains(err, "invalid request")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetWhiteList_ErrPageSize() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	_, err := k.GetWhitelist(ctx, &types.WhitelistRequest{
		Pagination: &query.PageRequest{
			Limit: 12000,
		},
	})

	suite.Require().ErrorContains(err, "page size greater than max")
}

func (suite *PerpetualKeeperTestSuite) TestQueryGetWhiteList_Successful() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx

	k.WhitelistAddress(ctx, sdktypes.AccAddress(sample.AccAddress()))
	_, err := k.GetWhitelist(ctx, &types.WhitelistRequest{
		Pagination: &query.PageRequest{
			Limit: 90,
		},
	})

	suite.Require().Nil(err)
}
