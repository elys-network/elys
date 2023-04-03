package types_test

import (
	"testing"

	"github.com/elys-network/elys/x/oracle/types"
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
				PortId: types.PortID,
				AssetInfoList: []types.AssetInfo{
					{
						Denom: "satoshi",
					},
					{
						Denom: "wei",
					},
				},
				PriceList: []types.Price{
					{
						Asset: "asset0",
					},
					{
						Asset: "asset1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated assetInfo",
			genState: &types.GenesisState{
				AssetInfoList: []types.AssetInfo{
					{
						Denom: "satoshi",
					},
					{
						Denom: "satoshi",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated price",
			genState: &types.GenesisState{
				PriceList: []types.Price{
					{
						Asset: "asset0",
					},
					{
						Asset: "asset0",
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
