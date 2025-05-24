package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/burner/types"
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
				HistoryList: []types.History{
					{
						Block:       1,
						BurnedCoins: sdk.Coins{sdk.NewInt64Coin("uusdc", 1)},
					},
					{
						Block:       2,
						BurnedCoins: sdk.Coins{sdk.NewInt64Coin("uatom", 1)},
					},
				},
				Params: types.DefaultParams(),
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated history",
			genState: &types.GenesisState{
				HistoryList: []types.History{
					{
						Block:       1,
						BurnedCoins: sdk.Coins{sdk.NewInt64Coin("uusdc", 1)},
					},
					{
						Block:       1,
						BurnedCoins: sdk.Coins{sdk.NewInt64Coin("uusdc", 1)},
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
