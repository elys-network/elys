package keeper_test

import (
	"errors"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/v7/testutil/keeper"
	"github.com/elys-network/elys/v7/x/tokenomics/keeper"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
)

func TestGenesisInflationMsgServerUpdate(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	inflation := &types.InflationEntry{
		LmRewards:         10,
		IcsStakingRewards: 10,
		CommunityFund:     10,
		StrategicReserve:  10,
		TeamTokensVested:  10,
	}

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateGenesisInflation
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateGenesisInflation{Authority: authority,
				Inflation:             inflation,
				SeedVesting:           10,
				StrategicSalesVesting: 5,
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateGenesisInflation{Authority: "B",
				Inflation:             inflation,
				SeedVesting:           10,
				StrategicSalesVesting: 5,
			},
			err: errors.New("invalid authority"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			expected := &types.MsgUpdateGenesisInflation{Authority: authority,
				Inflation:             inflation,
				SeedVesting:           10,
				StrategicSalesVesting: 8,
			}
			_, err := srv.UpdateGenesisInflation(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateGenesisInflation(ctx, tc.request)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				rst, found := k.GetGenesisInflation(ctx)
				require.True(t, found)
				require.Equal(t, expected.Authority, rst.Authority)
			}
		})
	}
}
