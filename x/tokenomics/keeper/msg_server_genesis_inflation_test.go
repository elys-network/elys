package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
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
			err: sdkerrors.ErrInvalidAddress,
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
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetGenesisInflation(ctx)
				require.True(t, found)
				require.Equal(t, expected.Authority, rst.Authority)
			}
		})
	}
}
