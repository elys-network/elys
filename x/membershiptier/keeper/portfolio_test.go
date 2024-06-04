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
		items[i].MinimumToday = sdk.NewDec(1000)

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
			ctx.BlockTime().Format("2020-02-01"),
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
		val := keeper.GetPortfolio(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
			ctx.BlockTime().Format("2020-02-01"),
		)
		require.EqualValues(t, []types.Portfolio([]types.Portfolio(nil)), val)
	}
}

func TestPortfolioGetMinimumToday(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePortfolio(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
			item.Token.Denom,
		)
		_, found := keeper.GetPortfolioMinimumToday(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
			ctx.BlockTime().Format("2020-02-01"),
			item.Token.Denom,
		)
		require.False(t, found)
	}
}

func TestProcessPortfolioChange(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)

	keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, "cosmos10t3g865e53yhhzvwwr5ldg50yq7vdwwfemrdeg", "2", sdk.NewInt(1000))

	total_port, tier, disc := keeper.GetMembershipTier(ctx, "cosmos10t3g865e53yhhzvwwr5ldg50yq7vdwwfemrdeg")
	require.Equal(t, sdk.NewDec(100), total_port)
	require.Equal(t, "bronze", tier)
	require.Equal(t, 0, disc)
}

func TestPortfolioGetAll(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPortfolio(ctx)),
	)
}
