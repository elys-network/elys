package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/testutil/nullify"
	"github.com/elys-network/elys/v4/x/burner/keeper"
	"github.com/elys-network/elys/v4/x/burner/types"
	"github.com/stretchr/testify/require"
)

func createNHistory(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.History {
	items := make([]types.History, n)
	for i := range items {
		items[i].Block = uint64(i)
		items[i].BurnedCoins = sdk.NewCoins(sdk.NewInt64Coin("uatom", int64(i)))

		keeper.SetHistory(ctx, items[i])
	}
	return items
}

func TestHistoryGet(t *testing.T) {
	keeper, ctx, _ := keepertest.BurnerKeeper(t)
	items := createNHistory(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetHistory(ctx, item.Block)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestHistoryRemove(t *testing.T) {
	keeper, ctx, _ := keepertest.BurnerKeeper(t)
	items := createNHistory(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveHistory(ctx, item.Block)
		_, found := keeper.GetHistory(ctx, item.Block)
		require.False(t, found)
	}
}

func TestHistoryGetAll(t *testing.T) {
	keeper, ctx, _ := keepertest.BurnerKeeper(t)
	items := createNHistory(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllHistory(ctx)),
	)
}
