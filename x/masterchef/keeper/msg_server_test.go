package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/masterchef/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MasterchefKeeper(t)

	k.SetPoolInfo(ctx, types.PoolInfo{
		PoolId:               2,
		RewardWallet:         "elys1d96rzrky937s3s397g5xh5qvcwgkeqysmh8sg2kn359fhfvzeyrsnalu2u",
		Multiplier:           sdk.MustNewDecFromStr("1.00"),
		GasApr:               sdk.MustNewDecFromStr("0.00"),
		EdenApr:              sdk.MustNewDecFromStr("0.50"),
		DexApr:               sdk.MustNewDecFromStr("0.00"),
		ExternalIncentiveApr: sdk.MustNewDecFromStr("0.00"),
		EnableEdenRewards:    false,
	})

	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
