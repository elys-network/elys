package keeper_test

import (
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/membershiptier/keeper"
	"github.com/elys-network/elys/x/membershiptier/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

const userAddress string = "cosmos10t3g865e53yhhzvwwr5ldg50yq7vdwwfemrdeg"

func createNPortfolio(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Portfolio {
	items := make([]types.Portfolio, n)
	for i := range items {
		items[i].Creator = strconv.Itoa(i)
		items[i].MinimumToday = sdk.NewDec(1000)
		items[i].Denom = strconv.Itoa(i)
		items[i].Assetkey = types.LiquidKeyPrefix
		items[i].MinimumToday = sdk.NewDec(100)
		items[i].Amount = 100

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
			keeper.GetDateFromBlock(ctx.BlockTime()),
		)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst[0]),
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
			keeper.GetDateFromBlock(ctx.BlockTime()),
		)
		val := keeper.GetPortfolio(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
			keeper.GetDateFromBlock(ctx.BlockTime()),
		)
		require.EqualValues(t, []types.Portfolio([]types.Portfolio(nil)), val)
	}
}

func TestPortfolioGetMinimumToday(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	for _, item := range items {
		_, found := keeper.GetPortfolioMinimumToday(ctx,
			item.Creator,
			types.LiquidKeyPrefix,
			keeper.GetDateFromBlock(ctx.BlockTime()),
			item.Denom,
		)
		require.True(t, found)
	}
}

func TestProcessPortfolioChange(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(1000))

	total_port, tier, disc := keeper.GetMembershipTier(ctx, userAddress)
	require.Equal(t, sdk.NewDec(0), total_port)
	require.Equal(t, "bronze", tier)
	require.Equal(t, uint64(0), disc)

	for i := 0; i < 9; i++ {
		keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(1000))
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 24))
	}
	keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(1000))

	total_port, tier, disc = keeper.GetMembershipTier(ctx, userAddress)
	require.Equal(t, sdk.NewDec(1000), total_port)
	require.Equal(t, "bronze", tier)
	require.Equal(t, uint64(0), disc)

	for i := 0; i < 9; i++ {
		keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(100000))
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 24))
	}
	keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(400000))

	total_port, tier, disc = keeper.GetMembershipTier(ctx, userAddress)
	require.Equal(t, sdk.NewDec(100000), total_port)
	require.Equal(t, "silver", tier)
	require.Equal(t, uint64(10), disc)

	for i := 0; i < 9; i++ {
		keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(500000))
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 24))
	}
	keeper.ProcessPortfolioChange(ctx, types.PerpetualKeyPrefix, userAddress, "2", sdk.NewInt(500000))

	total_port, tier, disc = keeper.GetMembershipTier(ctx, userAddress)
	require.Equal(t, sdk.NewDec(500000), total_port)
	require.Equal(t, "platinum", tier)
	require.Equal(t, uint64(30), disc)
}

func TestPortfolioGetAll(t *testing.T) {
	keeper, ctx := keepertest.MembershiptierKeeper(t)
	items := createNPortfolio(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPortfolio(ctx, keeper.GetDateFromBlock(ctx.BlockTime()))),
	)
}
