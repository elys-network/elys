package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/v5/testutil/keeper"
	"github.com/elys-network/elys/v5/x/burner/keeper"
	"github.com/elys-network/elys/v5/x/burner/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx, _ := keepertest.BurnerKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
