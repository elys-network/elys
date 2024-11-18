package keeper_test

import (
	"context"
	"testing"

	sdkmath "cosmossdk.io/math"

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
		Multiplier:           sdkmath.LegacyMustNewDecFromStr("1.00"),
		GasApr:               sdkmath.LegacyMustNewDecFromStr("0.00"),
		EdenApr:              sdkmath.LegacyMustNewDecFromStr("0.50"),
		DexApr:               sdkmath.LegacyMustNewDecFromStr("0.00"),
		ExternalIncentiveApr: sdkmath.LegacyMustNewDecFromStr("0.00"),
		EnableEdenRewards:    false,
	})

	return keeper.NewMsgServerImpl(*k), ctx
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
