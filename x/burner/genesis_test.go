package burner_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/burner"
	"github.com/elys-network/elys/x/burner/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		HistoryList: []types.History{
			{
				Timestamp: "0",
				Denom:     "0",
			},
			{
				Timestamp: "1",
				Denom:     "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BurnerKeeper(t)
	burner.InitGenesis(ctx, *k, genesisState)
	got := burner.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.HistoryList, got.HistoryList)
	// this line is used by starport scaffolding # genesis/test/assert
}
