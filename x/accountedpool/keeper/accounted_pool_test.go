package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/accountedpool/keeper"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNAccountedPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.AccountedPool {
	items := make([]types.AccountedPool, n)
	for i := range items {
		items[i].PoolId = (uint64)(i)
		items[i].TotalShares = sdk.NewCoin("lpshare", sdk.ZeroInt())
		items[i].PoolAssets = []ammtypes.PoolAsset{}
		items[i].TotalWeight = sdk.ZeroInt()

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
