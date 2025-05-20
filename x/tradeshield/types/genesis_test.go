package types_test

import (
	"testing"

	"github.com/elys-network/elys/v4/x/tradeshield/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	params := types.DefaultParams()
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
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: params,
				PendingSpotOrderList: []types.SpotOrder{
					{
						OrderId: 0,
					},
					{
						OrderId: 1,
					},
				},
				PendingSpotOrderCount: 2,
				PendingPerpetualOrderList: []types.PerpetualOrder{
					{
						OrderId: 0,
					},
					{
						OrderId: 1,
					},
				},
				PendingPerpetualOrderCount: 2,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated pendingSpotOrder",
			genState: &types.GenesisState{
				PendingSpotOrderList: []types.SpotOrder{
					{
						OrderId: 0,
					},
					{
						OrderId: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid pendingSpotOrder count",
			genState: &types.GenesisState{
				PendingSpotOrderList: []types.SpotOrder{
					{
						OrderId: 1,
					},
				},
				PendingSpotOrderCount: 0,
			},
			valid: false,
		},
		{
			desc: "duplicated pendingPerpetualOrder",
			genState: &types.GenesisState{
				PendingPerpetualOrderList: []types.PerpetualOrder{
					{
						OrderId: 0,
					},
					{
						OrderId: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid pendingPerpetualOrder count",
			genState: &types.GenesisState{
				PendingPerpetualOrderList: []types.PerpetualOrder{
					{
						OrderId: 1,
					},
				},
				PendingPerpetualOrderCount: 0,
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
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
