package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/tokenomics/keeper"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func createNAirdrop(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Airdrop {
	items := make([]types.Airdrop, n)
	for i := range items {
		items[i].Intent = strconv.Itoa(i)

		keeper.SetAirdrop(ctx, items[i])
	}
	return items
}

func TestAirdropGet(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	items := createNAirdrop(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetAirdrop(ctx,
			item.Intent,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestAirdropRemove(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	items := createNAirdrop(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAirdrop(ctx,
			item.Intent,
		)
		_, found := keeper.GetAirdrop(ctx,
			item.Intent,
		)
		require.False(t, found)
	}
}

func TestAirdropGetAll(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	items := createNAirdrop(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAirdrop(ctx)),
	)
}
