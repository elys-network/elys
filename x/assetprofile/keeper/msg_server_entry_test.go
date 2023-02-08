package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/assetprofile/keeper"
	"github.com/elys-network/elys/x/assetprofile/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestEntryMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.AssetprofileKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateEntry{Authority: authority,
			BaseDenom: strconv.Itoa(i),
		}
		_, err := srv.CreateEntry(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetEntry(ctx,
			expected.BaseDenom,
		)
		require.True(t, found)
		require.Equal(t, expected.Authority, rst.Authority)
	}
}

func TestEntryMsgServerUpdate(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateEntry
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateEntry{Authority: authority,
				BaseDenom: strconv.Itoa(0),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateEntry{Authority: "B",
				BaseDenom: strconv.Itoa(0),
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateEntry{Authority: authority,
				BaseDenom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.AssetprofileKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateEntry{Authority: authority,
				BaseDenom: strconv.Itoa(0),
			}
			_, err := srv.CreateEntry(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateEntry(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetEntry(ctx,
					expected.BaseDenom,
				)
				require.True(t, found)
				require.Equal(t, expected.Authority, rst.Authority)
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
			request: &types.MsgDeleteEntry{Authority: authority,
				BaseDenom: strconv.Itoa(0),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgDeleteEntry{Authority: "B",
				BaseDenom: strconv.Itoa(0),
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteEntry{Authority: authority,
				BaseDenom: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.AssetprofileKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateEntry(wctx, &types.MsgCreateEntry{Authority: authority,
				BaseDenom: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteEntry(wctx, tc.request)
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
