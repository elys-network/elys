package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/membershiptier/keeper"
	"github.com/elys-network/elys/x/membershiptier/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPortfolio(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Portfolio {
	items := make([]types.Portfolio, n)
	for i := range items {
		items[i].Creator = strconv.Itoa(i)

		keeper.SetPortfolio(ctx, items[i], types.LiquidKeyPrefix)
	}
	return items
}

func TestPortfolioGet(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		rst := keeper.GetPortfolio(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
		)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestPortfolioRemove(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePortfolio(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
			item.Token.Denom,
		)
		keeper.GetPortfolio(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
		)
		// require.
		// require.False(t, found)
	}
}

func TestPortfolioGetAll(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPortfolio(ctx)),
	)
}
