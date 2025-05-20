package keeper_test

import (
	"errors"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/x/tokenomics/keeper"
	"github.com/elys-network/elys/v4/x/tokenomics/types"
)

func TestTimeBasedInflationMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenomicsKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	description := "test"
	inflation := &types.InflationEntry{
		LmRewards:         10,
		IcsStakingRewards: 10,
		CommunityFund:     10,
		StrategicReserve:  10,
		TeamTokensVested:  10,
	}

	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateTimeBasedInflation{
			Authority:        authority,
			StartBlockHeight: uint64(i),
			EndBlockHeight:   100,
			Description:      description,
			Inflation:        inflation,
		}
		_, err := srv.CreateTimeBasedInflation(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetTimeBasedInflation(ctx,
			expected.StartBlockHeight,
			expected.EndBlockHeight,
		)
		require.True(t, found)
		require.Equal(t, expected.Authority, rst.Authority)
	}
}

func TestTimeBasedInflationMsgServerUpdate(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	description := "test"
	inflation := &types.InflationEntry{
		LmRewards:         10,
		IcsStakingRewards: 10,
		CommunityFund:     10,
		StrategicReserve:  10,
		TeamTokensVested:  10,
	}

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateTimeBasedInflation
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 100,
				EndBlockHeight:   1000,
				Description:      description,
				Inflation:        inflation,
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateTimeBasedInflation{
				Authority:        "B",
				StartBlockHeight: 0,
				EndBlockHeight:   0,
				Description:      description,
				Inflation:        inflation,
			},
			err: errors.New("invalid authority"),
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 100000,
				EndBlockHeight:   100001,
				Description:      description,
				Inflation:        inflation,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)

			expected := &types.MsgCreateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 100,
				EndBlockHeight:   1000,
				Description:      description,
				Inflation:        inflation,
			}
			_, err := srv.CreateTimeBasedInflation(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateTimeBasedInflation(ctx, tc.request)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				rst, found := k.GetTimeBasedInflation(ctx,
					expected.StartBlockHeight,
					expected.EndBlockHeight,
				)
				require.True(t, found)
				require.Equal(t, expected.Authority, rst.Authority)
			}
		})
	}
}

func TestTimeBasedInflationMsgServerDelete(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteTimeBasedInflation
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 10,
				EndBlockHeight:   100,
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgDeleteTimeBasedInflation{
				Authority:        "B",
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
			err: errors.New("invalid authority"),
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 100000,
				EndBlockHeight:   120000,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)

			_, err := srv.CreateTimeBasedInflation(ctx, &types.MsgCreateTimeBasedInflation{
				Description:      "Test create time based inflation",
				Authority:        authority,
				StartBlockHeight: 10,
				EndBlockHeight:   100,
				Inflation: &types.InflationEntry{
					LmRewards:         10,
					IcsStakingRewards: 10,
					CommunityFund:     10,
					StrategicReserve:  10,
					TeamTokensVested:  10,
				},
			})
			require.NoError(t, err)
			_, err = srv.DeleteTimeBasedInflation(ctx, tc.request)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				_, found := k.GetTimeBasedInflation(ctx,
					tc.request.StartBlockHeight,
					tc.request.EndBlockHeight,
				)
				require.False(t, found)
			}
		})
	}
}
