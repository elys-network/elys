package keeper_test

import (
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryGetConsolidatedPriceInvalidRequest() {
	_, err := suite.app.TierKeeper.GetConsolidatedPrice(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryGetConsolidatedPriceSuccessful() {
	_, err := suite.app.TierKeeper.GetConsolidatedPrice(suite.ctx, &types.QueryGetConsolidatedPriceRequest{
		Denom: ptypes.ATOM,
	})

	suite.Require().NoError(err)

	_, err = suite.app.TierKeeper.GetConsolidatedPrice(suite.ctx, &types.QueryGetConsolidatedPriceRequest{
		Denom: ptypes.Eden,
	})

	suite.Require().NoError(err)

}
