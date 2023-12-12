package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func TestTimeBasedInflationMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenomicsKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateTimeBasedInflation{
			Authority:        authority,
			StartBlockHeight: uint64(i),
			EndBlockHeight:   uint64(i),
		}
		_, err := srv.CreateTimeBasedInflation(wctx, expected)
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

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateTimeBasedInflation
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateTimeBasedInflation{
				Authority:        "B",
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 100000,
				EndBlockHeight:   100000,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			}
			_, err := srv.CreateTimeBasedInflation(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateTimeBasedInflation(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
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
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgDeleteTimeBasedInflation{
				Authority:        "B",
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 100000,
				EndBlockHeight:   100000,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateTimeBasedInflation(wctx, &types.MsgCreateTimeBasedInflation{
				Authority:        authority,
				StartBlockHeight: 0,
				EndBlockHeight:   0,
			})
			require.NoError(t, err)
			_, err = srv.DeleteTimeBasedInflation(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
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
