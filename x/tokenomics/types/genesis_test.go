package types_test

import (
	"testing"

	"github.com/elys-network/elys/v7/x/tokenomics/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
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
					SeedVesting:           94,
					StrategicSalesVesting: 51,
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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated airdrop",
			genState: &types.GenesisState{
				AirdropList: []types.Airdrop{
					{
						Intent: "0",
					},
					{
						Intent: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated timeBasedInflation",
			genState: &types.GenesisState{
				TimeBasedInflationList: []types.TimeBasedInflation{
					{
						StartBlockHeight: 0,
						EndBlockHeight:   0,
					},
					{
						StartBlockHeight: 0,
						EndBlockHeight:   0,
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
