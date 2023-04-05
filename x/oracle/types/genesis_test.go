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
				Params: types.DefaultParams(),
				AssetInfos: []types.AssetInfo{
					{
						Denom: "satoshi",
					},
					{
						Denom: "wei",
					},
				},
				Prices: []types.Price{
					{
						Asset: "asset0",
					},
					{
						Asset: "asset1",
					},
				},
				PriceFeeders: []types.PriceFeeder{
					{
						Feeder: "0",
					},
					{
						Feeder: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated assetInfo",
			genState: &types.GenesisState{
				AssetInfos: []types.AssetInfo{
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
				Prices: []types.Price{
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
		{
			desc: "duplicated priceFeeder",
			genState: &types.GenesisState{
				PriceFeeders: []types.PriceFeeder{
					{
						Feeder: "0",
					},
					{
						Feeder: "0",
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
