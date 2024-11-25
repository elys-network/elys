package keeper_test

import (
	"context"

	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) setupMsgServer() (types.MsgServer, context.Context) {
	return keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper), suite.ctx
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServer() {
	ms, ctx := suite.setupMsgServer()
	suite.Require().NotNil(ms)
	suite.Require().NotNil(ctx)
}
