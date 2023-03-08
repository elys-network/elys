package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func TestGenesisInflationMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenomicsKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	authority := "A"
	expected := &types.MsgCreateGenesisInflation{Authority: authority}
	_, err := srv.CreateGenesisInflation(wctx, expected)
	require.NoError(t, err)
	rst, found := k.GetGenesisInflation(ctx)
	require.True(t, found)
	require.Equal(t, expected.Authority, rst.Authority)
}

func TestGenesisInflationMsgServerUpdate(t *testing.T) {
	authority := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateGenesisInflation
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateGenesisInflation{Authority: authority},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateGenesisInflation{Authority: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateGenesisInflation{Authority: authority}
			_, err := srv.CreateGenesisInflation(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateGenesisInflation(wctx, tc.request)
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

func TestGenesisInflationMsgServerDelete(t *testing.T) {
	authority := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteGenesisInflation
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteGenesisInflation{Authority: authority},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteGenesisInflation{Authority: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateGenesisInflation(wctx, &types.MsgCreateGenesisInflation{Authority: authority})
			require.NoError(t, err)
			_, err = srv.DeleteGenesisInflation(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetGenesisInflation(ctx)
				require.False(t, found)
			}
		})
	}
}
