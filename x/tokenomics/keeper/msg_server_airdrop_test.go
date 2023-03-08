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
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestAirdropMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenomicsKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateAirdrop{Authority: authority,
			Intent: strconv.Itoa(i),
		}
		_, err := srv.CreateAirdrop(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetAirdrop(ctx,
			expected.Intent,
		)
		require.True(t, found)
		require.Equal(t, expected.Authority, rst.Authority)
	}
}

func TestAirdropMsgServerUpdate(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateAirdrop
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateAirdrop{Authority: authority,
				Intent: strconv.Itoa(0),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateAirdrop{Authority: "B",
				Intent: strconv.Itoa(0),
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateAirdrop{Authority: authority,
				Intent: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateAirdrop{Authority: authority,
				Intent: strconv.Itoa(0),
			}
			_, err := srv.CreateAirdrop(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateAirdrop(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetAirdrop(ctx,
					expected.Intent,
				)
				require.True(t, found)
				require.Equal(t, expected.Authority, rst.Authority)
			}
		})
	}
}

func TestAirdropMsgServerDelete(t *testing.T) {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteAirdrop
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteAirdrop{Authority: authority,
				Intent: strconv.Itoa(0),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgDeleteAirdrop{Authority: "B",
				Intent: strconv.Itoa(0),
			},
			err: govtypes.ErrInvalidSigner,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteAirdrop{Authority: authority,
				Intent: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateAirdrop(wctx, &types.MsgCreateAirdrop{Authority: authority,
				Intent: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteAirdrop(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetAirdrop(ctx,
					tc.request.Intent,
				)
				require.False(t, found)
			}
		})
	}
}
