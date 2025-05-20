package keeper_test

import (
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	keepertest "github.com/elys-network/elys/v4/testutil/keeper"
	"github.com/elys-network/elys/v4/x/assetprofile/keeper"
	"github.com/elys-network/elys/v4/x/assetprofile/types"
)

func TestEntryMsgServerUpdate(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateEntry
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateEntry{
				Authority: authority,
				BaseDenom: strconv.Itoa(0),
				Decimals:  6,
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateEntry{
				Authority: "B",
				BaseDenom: strconv.Itoa(0),
				Decimals:  6,
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateEntry{
				Authority: authority,
				BaseDenom: strconv.Itoa(100000),
				Decimals:  6,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.AssetprofileKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			expected := &types.MsgAddEntry{
				Creator:   authority,
				BaseDenom: strconv.Itoa(0),
				Decimals:  6,
			}
			_, err := srv.AddEntry(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateEntry(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetEntry(ctx,
					expected.BaseDenom,
				)
				require.True(t, found)
				//require.Equal(t, expected.Authority, rst.Authority)
			}
		})
	}
}

func TestEntryMsgServerDelete(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteEntry
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteEntry{
				Authority: authority,
				BaseDenom: strconv.Itoa(0),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgDeleteEntry{
				Authority: "B",
				BaseDenom: strconv.Itoa(0),
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteEntry{
				Authority: authority,
				BaseDenom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.AssetprofileKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)

			_, err := srv.AddEntry(ctx, &types.MsgAddEntry{
				Creator:   authority,
				Decimals:  6,
				BaseDenom: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteEntry(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetEntry(ctx,
					tc.request.BaseDenom,
				)
				require.False(t, found)
			}
		})
	}
}
