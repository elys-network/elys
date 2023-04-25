package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/incentive/keeper"
	"github.com/elys-network/elys/x/incentive/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNElysDelegator(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ElysDelegator {
	items := make([]types.ElysDelegator, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetElysDelegator(ctx, items[i])
	}
	return items
}

func TestElysDelegatorGet(t *testing.T) {
	keeper, ctx := keepertest.IncentiveKeeper(t)
	items := createNElysDelegator(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetElysDelegator(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestElysDelegatorRemove(t *testing.T) {
	keeper, ctx := keepertest.IncentiveKeeper(t)
	items := createNElysDelegator(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveElysDelegator(ctx,
			item.Index,
		)
		_, found := keeper.GetElysDelegator(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestElysDelegatorGetAll(t *testing.T) {
	keeper, ctx := keepertest.IncentiveKeeper(t)
	items := createNElysDelegator(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllElysDelegator(ctx)),
	)
}
