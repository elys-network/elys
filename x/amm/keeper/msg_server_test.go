package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx, _, _ := keepertest.AmmKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
