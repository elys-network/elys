package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/require"
)

func createNPendingSpotOrder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PendingSpotOrder {
	items := make([]types.PendingSpotOrder, n)
	for i := range items {
		items[i].Id = keeper.AppendPendingSpotOrder(ctx, items[i])
	}
	return items
}

func TestPendingSpotOrderGet(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPendingSpotOrder(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPendingSpotOrderRemove(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePendingSpotOrder(ctx, item.Id)
		_, found := keeper.GetPendingSpotOrder(ctx, item.Id)
		require.False(t, found)
	}
}

func TestPendingSpotOrderGetAll(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPendingSpotOrder(ctx)),
	)
}

func TestPendingSpotOrderCount(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPendingSpotOrderCount(ctx))
}
