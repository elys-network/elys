package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/tokenomics/keeper"
	"github.com/elys-network/elys/v6/x/tokenomics/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.TokenomicsKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
