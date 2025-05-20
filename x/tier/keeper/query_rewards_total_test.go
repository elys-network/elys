package keeper_test

import (
	"github.com/elys-network/elys/v4/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryRewardsTotalRequest() {
	k := suite.app.TierKeeper
	_, err := k.RewardsTotal(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryRewardsTotalSuccessful() {
	k := suite.app.TierKeeper
	_, err := k.RewardsTotal(suite.ctx, &types.QueryRewardsTotalRequest{
		User: suite.account.String(),
	})

	suite.Require().NoError(err)
}
