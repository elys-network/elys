package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx, _, _, _, _ := keepertest.TradeshieldKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServer() {
	ms, ctx := setupMsgServer(suite.T())
	suite.Require().NotNil(ms)
	suite.Require().NotNil(ctx)
}
