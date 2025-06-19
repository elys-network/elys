package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/testutil/nullify"
	"github.com/elys-network/elys/v6/x/accountedpool/keeper"
	"github.com/elys-network/elys/v6/x/accountedpool/types"
	"github.com/stretchr/testify/require"
)

func createNAccountedPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.AccountedPool {
	items := make([]types.AccountedPool, n)
	for i := range items {
		items[i].PoolId = (uint64)(i)

		keeper.SetAccountedPool(ctx, items[i])
	}
	return items
}

func TestAccountedPoolGet(t *testing.T) {
	keeper, ctx := keepertest.AccountedPoolKeeper(t)
	items := createNAccountedPool(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetAccountedPool(ctx,
			item.PoolId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestAccountedPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.AccountedPoolKeeper(t)
	items := createNAccountedPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAccountedPool(ctx,
			item.PoolId,
		)
		_, found := keeper.GetAccountedPool(ctx,
			item.PoolId,
		)
		require.False(t, found)
	}
}

func TestAccountedPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.AccountedPoolKeeper(t)
	items := createNAccountedPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAccountedPool(ctx)),
	)
}
