package stablestake_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/testutil/nullify"
	"github.com/elys-network/elys/v4/x/stablestake"
	"github.com/elys-network/elys/v4/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.StablestakeKeeper(t)
	stablestake.InitGenesis(ctx, *k, genesisState)
	got := stablestake.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
