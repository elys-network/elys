package types_test

import (
	"testing"

	"github.com/elys-network/elys/v4/x/amm/types"
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
				PoolList: []types.Pool{
					{
						PoolId: 0,
					},
					{
						PoolId: 1,
					},
				},
				DenomLiquidityList: []types.DenomLiquidity{
					{
						Denom: "0",
					},
					{
						Denom: "1",
					},
				},
				Params: types.DefaultParams(),
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated pool",
			genState: &types.GenesisState{
				PoolList: []types.Pool{
					{
						PoolId: 0,
					},
					{
						PoolId: 0,
					},
				},
				Params: types.DefaultParams(),
			},
			valid: false,
		},
		{
			desc: "duplicated denomLiquidity",
			genState: &types.GenesisState{
				DenomLiquidityList: []types.DenomLiquidity{
					{
						Denom: "0",
					},
					{
						Denom: "0",
					},
				},
				Params: types.DefaultParams(),
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
