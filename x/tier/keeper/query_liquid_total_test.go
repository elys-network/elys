package keeper_test

import (
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryLiquidTotalInvalidRequest() {
	k := suite.app.TierKeeper
	_, err := k.LiquidTotal(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryLiquidTotalSuccessful() {
	k := suite.app.TierKeeper

	_, err := k.LiquidTotal(suite.ctx, &types.QueryLiquidTotalRequest{
		User: suite.account.String(),
	})

	suite.Require().NoError(err)
}
