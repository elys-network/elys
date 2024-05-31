package membershiptier_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/membershiptier"
	"github.com/elys-network/elys/x/membershiptier/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PortfolioList: []types.Portfolio{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MembershiptierKeeper(t)
	membershiptier.InitGenesis(ctx, *k, genesisState)
	got := membershiptier.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PortfolioList, got.PortfolioList)
	// this line is used by starport scaffolding # genesis/test/assert
}
