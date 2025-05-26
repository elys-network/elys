package keeper_test

import (
	"cosmossdk.io/math"
	stablestakekeeper "github.com/elys-network/elys/v5/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/v5/x/stablestake/types"
	"github.com/elys-network/elys/v5/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *TierKeeperTestSuite) TestQueryStakedInvalidRequest() {
	k := suite.app.TierKeeper
	_, err := k.Staked(suite.ctx, nil)

	want := status.Error(codes.InvalidArgument, "invalid request")

	suite.Require().ErrorIs(err, want)
}

func (suite *TierKeeperTestSuite) TestQueryStaked() {

	sender := suite.account

	msgServer := stablestakekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)

	//STAKE USDC
	_, err := msgServer.Bond(suite.ctx, &stablestaketypes.MsgBond{
		Creator: sender.String(),
		Amount:  math.NewInt(100000000),
		PoolId:  1,
	})
	suite.Require().NoError(err)

	//TESTING STAKED FUNCTION.
	k := suite.app.TierKeeper
	r, err := k.Staked(suite.ctx, &types.QueryStakedRequest{
		User: sender.String(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(r.Commitments.TruncateInt(), math.NewInt(100))
}
