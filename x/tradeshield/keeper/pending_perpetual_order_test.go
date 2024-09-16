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

func createNPendingPerpetualOrder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PerpetualOrder {
	items := make([]types.PerpetualOrder, n)
	for i := range items {
		items[i].OrderId = keeper.AppendPendingPerpetualOrder(ctx, items[i])
	}
	return items
}

func TestPendingPerpetualOrderGet(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPendingPerpetualOrder(ctx, item.OrderId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPendingPerpetualOrderRemove(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePendingPerpetualOrder(ctx, item.OrderId)
		_, found := keeper.GetPendingPerpetualOrder(ctx, item.OrderId)
		require.False(t, found)
	}
}

func TestPendingPerpetualOrderGetAll(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPendingPerpetualOrder(ctx)),
	)
}

func TestPendingPerpetualOrderCount(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPendingPerpetualOrderCount(ctx))
}
