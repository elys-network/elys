package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/tokenomics/keeper"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
)

func createTestGenesisInflation(keeper *keeper.Keeper, ctx sdk.Context) types.GenesisInflation {
	item := types.GenesisInflation{}
	keeper.SetGenesisInflation(ctx, item)
	return item
}

func TestGenesisInflationGet(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	item := createTestGenesisInflation(keeper, ctx)
	rst, found := keeper.GetGenesisInflation(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestGenesisInflationRemove(t *testing.T) {
	keeper, ctx := keepertest.TokenomicsKeeper(t)
	createTestGenesisInflation(keeper, ctx)
	keeper.RemoveGenesisInflation(ctx)
	_, found := keeper.GetGenesisInflation(ctx)
	require.False(t, found)
}
