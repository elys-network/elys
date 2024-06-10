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

//const userAddress string = "cosmos10t3g865e53yhhzvwwr5ldg50yq7vdwwfemrdeg"

func createNPortfolio(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Portfolio {
	items := make([]types.Portfolio, n)
	for i := range items {
		items[i].Creator = strconv.Itoa(i)
		items[i].MinimumToday = sdk.NewDec(1000)
		items[i].Denom = strconv.Itoa(i)
		items[i].Assetkey = types.LiquidKeyPrefix
		items[i].MinimumToday = sdk.NewDec(100)
		items[i].Amount = 100

		keeper.SetPortfolio(ctx, keeper.GetDateFromBlock(ctx.BlockTime()), items[i].Creator, items[i])
	}
	return items
}

// TODO
// Test case 1: Only native balances assets
// 2: native + amm pool token
// 3: rewards
// 4: native + perpetual
func TestPortfolioGet(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		_, found := keeper.GetPortfolio(ctx,
			item.Creator,
			keeper.GetDateFromBlock(ctx.BlockTime()),
		)
		require.True(t, found)
	}
}
func TestPortfolioRemove(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePortfolio(ctx,
			item.Creator,
			keeper.GetDateFromBlock(ctx.BlockTime()),
		)
		_, found := keeper.GetPortfolio(ctx,
			item.Creator,
			keeper.GetDateFromBlock(ctx.BlockTime()),
		)
		require.False(t, found)
	}
}

func TestPortfolioGetAll(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPortfolio(ctx, keeper.GetDateFromBlock(ctx.BlockTime()))),
	)
}
