package keeper_test

import (
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryCalculateDiscountInvalidRequest() {
	_, err := suite.app.TierKeeper.CalculateDiscount(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryCalculateInvalidAddress() {
	_, err := suite.app.TierKeeper.CalculateDiscount(suite.ctx, &types.QueryCalculateDiscountRequest{
		User: "invalidAddress",
	})

	suite.Require().Error(err)
}

func (suite *TierKeeperTestSuite) TestQueryCalculateDiscountSuccessful() {
	_, err := suite.app.TierKeeper.CalculateDiscount(suite.ctx, &types.QueryCalculateDiscountRequest{
		User: suite.account.String(),
	})

	suite.Require().NoError(err)
}
