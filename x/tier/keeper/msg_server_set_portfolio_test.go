package keeper_test

import (
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/elys-network/elys/v5/x/tier/keeper"
	"github.com/elys-network/elys/v5/x/tier/types"
)

func (suite *TierKeeperTestSuite) TestMsgSetPortfolio() {

	msgServer := keeper.NewMsgServerImpl(*suite.app.TierKeeper)
	_, err := msgServer.SetPortfolio(suite.ctx, &types.MsgSetPortfolio{
		User: sample.AccAddress(),
	})
	suite.Require().NoError(err)
}
