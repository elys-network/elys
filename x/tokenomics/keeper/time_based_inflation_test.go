package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/tokenomics/keeper"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func createNTimeBasedInflation(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.TimeBasedInflation {
	items := make([]types.TimeBasedInflation, n)
	for i := range items {
		items[i].StartBlockHeight = uint64(i)
		items[i].EndBlockHeight = uint64(i)

		keeper.SetTimeBasedInflation(ctx, items[i])
	}
	return items
}

func TestTimeBasedInflationGet(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	items := createNTimeBasedInflation(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTimeBasedInflation(ctx,
			item.StartBlockHeight,
			item.EndBlockHeight,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestTimeBasedInflationRemove(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	items := createNTimeBasedInflation(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTimeBasedInflation(ctx,
			item.StartBlockHeight,
			item.EndBlockHeight,
		)
		_, found := keeper.GetTimeBasedInflation(ctx,
			item.StartBlockHeight,
			item.EndBlockHeight,
		)
		require.False(t, found)
	}
}

func TestTimeBasedInflationGetAll(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	items := createNTimeBasedInflation(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllTimeBasedInflation(ctx)),
	)
}
