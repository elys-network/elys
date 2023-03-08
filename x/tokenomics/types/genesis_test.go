package types_test

import (
	"testing"

	"github.com/elys-network/elys/x/tokenomics/types"
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
					Inflation:             "92",
					SeedVesting:           94,
					StrategicSalesVesting: 51,
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
