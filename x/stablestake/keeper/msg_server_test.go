package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/stablestake/keeper"
	"github.com/elys-network/elys/v6/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.StablestakeKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
