package keeper_test

import (
	"github.com/elys-network/elys/v6/x/perpetual/keeper"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestMsgServerOpen() {
	k := suite.app.PerpetualKeeper
	ctx := suite.ctx
	msg := keeper.NewMsgServerImpl(*k)
	_, err := msg.Open(ctx, &types.MsgOpen{})
	suite.Require().Error(err)
}
