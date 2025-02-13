package types_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"testing"
	"time"

	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	currentTime := time.Now().Unix()
	tests := []struct {
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
			desc:     "empty genesis state",
			genState: &types.GenesisState{},
			valid:    false,
		},
		{
			desc: "invalid genesis state",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DebtList: []types.Debt{
					{
						Address: authtypes.NewModuleAddress("1").String(),
					},
					{
						Address: authtypes.NewModuleAddress("1").String(),
					},
				},
				InterestList: nil,
			},
			valid: false,
		},
		{
			desc: "invalid genesis state",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				InterestList: []types.InterestBlock{
					{
						BlockTime: currentTime,
					},
					{
						BlockTime: currentTime,
					},
				},
				DebtList: nil,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
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
