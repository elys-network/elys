package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/parameter/keeper"
	"github.com/elys-network/elys/v7/x/parameter/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ParameterKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
