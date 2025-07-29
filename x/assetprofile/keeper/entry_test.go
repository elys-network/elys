package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/assetprofile/keeper"
	"github.com/elys-network/elys/v7/x/assetprofile/types"
	"github.com/stretchr/testify/require"
)

func createNEntry(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Entry {
	items := make([]types.Entry, n)
	for i := range items {
		items[i].BaseDenom = strconv.Itoa(i)

		keeper.SetEntry(ctx, items[i])
	}
	return items
}

func TestEntryGet(t *testing.T) {
	keeper, ctx := keepertest.AssetprofileKeeper(t)
	items := createNEntry(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetEntry(ctx,
			item.BaseDenom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestEntryRemove(t *testing.T) {
	keeper, ctx := keepertest.AssetprofileKeeper(t)
	items := createNEntry(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveEntry(ctx,
			item.BaseDenom,
		)
		_, found := keeper.GetEntry(ctx,
			item.BaseDenom,
		)
		require.False(t, found)
	}
}

func TestEntryGetAll(t *testing.T) {
	keeper, ctx := keepertest.AssetprofileKeeper(t)
	items := createNEntry(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllEntry(ctx)),
	)
}
