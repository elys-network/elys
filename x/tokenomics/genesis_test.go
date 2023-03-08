package tokenomics_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tokenomics"
	"github.com/elys-network/elys/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AirdropList: []types.Airdrop{
			{
				Intent: "0",
			},
			{
				Intent: "1",
			},
		},
		GenesisInflation: &types.GenesisInflation{
			Inflation:             "46",
			SeedVesting:           85,
			StrategicSalesVesting: 5,
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TokenomicsKeeper(t)
	tokenomics.InitGenesis(ctx, *k, genesisState)
	got := tokenomics.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.AirdropList, got.AirdropList)
	require.Equal(t, genesisState.GenesisInflation, got.GenesisInflation)
	// this line is used by starport scaffolding # genesis/test/assert
}
