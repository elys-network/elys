package tokenomics_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/v5/testutil/keeper"
	"github.com/elys-network/elys/v5/testutil/nullify"
	"github.com/elys-network/elys/v5/x/tokenomics"
	"github.com/elys-network/elys/v5/x/tokenomics/types"
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
			Inflation: &types.InflationEntry{
				LmRewards:         10,
				IcsStakingRewards: 10,
				CommunityFund:     10,
				StrategicReserve:  10,
				TeamTokensVested:  10,
			},
			SeedVesting:           85,
			StrategicSalesVesting: 5,
		},
		TimeBasedInflationList: []types.TimeBasedInflation{
			{
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
			{
				StartBlockHeight: 1,
				EndBlockHeight:   1,
			},
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
	require.ElementsMatch(t, genesisState.TimeBasedInflationList, got.TimeBasedInflationList)
	// this line is used by starport scaffolding # genesis/test/assert
}
