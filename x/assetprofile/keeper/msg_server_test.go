package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/assetprofile/keeper"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AssetprofileKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
