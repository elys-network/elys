package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/burner/keeper"
	"github.com/elys-network/elys/x/burner/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNHistory(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.History {
	items := make([]types.History, n)
	for i := range items {
		items[i].Timestamp = strconv.Itoa(i)
		items[i].Denom = strconv.Itoa(i)

		keeper.SetHistory(ctx, items[i])
	}
	return items
}

func TestHistoryGet(t *testing.T) {
	keeper, ctx := keepertest.BurnerKeeper(t)
	items := createNHistory(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetHistory(ctx,
			item.Timestamp,
			item.Denom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestHistoryRemove(t *testing.T) {
	keeper, ctx := keepertest.BurnerKeeper(t)
	items := createNHistory(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveHistory(ctx,
			item.Timestamp,
			item.Denom,
		)
		_, found := keeper.GetHistory(ctx,
			item.Timestamp,
			item.Denom,
		)
		require.False(t, found)
	}
}

func TestHistoryGetAll(t *testing.T) {
	keeper, ctx := keepertest.BurnerKeeper(t)
	items := createNHistory(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllHistory(ctx)),
	)
}
