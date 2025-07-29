package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/assetprofile/keeper"
	"github.com/elys-network/elys/v7/x/assetprofile/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AssetprofileKeeper(t)
	return keeper.NewMsgServerImpl(*k), ctx
}
