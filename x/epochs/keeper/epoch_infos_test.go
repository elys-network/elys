package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/epochs/keeper"
	"github.com/elys-network/elys/x/epochs/types"
	"github.com/stretchr/testify/require"
)

func createEpochInfos(keeper *keeper.Keeper, ctx sdk.Context) []types.EpochInfo {

	items := make([]types.EpochInfo, 3)
	items[0] = types.EpochInfo{
		Identifier:            "daily",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	items[1] = types.EpochInfo{
		Identifier:            "weekly",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24 * 7,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	items[2] = types.EpochInfo{
		Identifier:            "monthly",
		StartTime:             time.Time{},
		Duration:              time.Hour * 24 * 30,
		CurrentEpoch:          0,
		CurrentEpochStartTime: time.Time{},
		EpochCountingStarted:  false,
	}
	for i := range items {
		keeper.SetEpochInfo(ctx, items[i])
	}
	return items
}

func TestEpochInfoGet(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	items := createEpochInfos(keeper, ctx)
	for _, item := range items {
		rst, found := keeper.GetEpochInfo(ctx,
			item.Identifier,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestEpochInfoRemove(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	items := createEpochInfos(keeper, ctx)
	for _, item := range items {
		keeper.DeleteEpochInfo(ctx,
			item.Identifier,
		)
		_, found := keeper.GetEpochInfo(ctx,
			item.Identifier,
		)
		require.False(t, found)
	}
}

func TestEntryGetAll(t *testing.T) {
	keeper, ctx := keepertest.EpochsKeeper(t)
	items := createEpochInfos(keeper, ctx)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.AllEpochInfos(ctx)),
	)
}
