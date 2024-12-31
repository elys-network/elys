package keeper_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/tokenomics/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func TestAirdropMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.TokenomicsKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateAirdrop{
			Authority: authority,
			Intent:    strconv.Itoa(i),
		}
		_, err := srv.CreateAirdrop(ctx, expected)
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
			request: &types.MsgUpdateAirdrop{
				Authority: authority,
				Intent:    strconv.Itoa(0),
				Amount:    200,
				Expiry:    uint64(time.Now().Unix()),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgUpdateAirdrop{
				Authority: "B",
				Intent:    strconv.Itoa(0),
				Amount:    200,
				Expiry:    uint64(time.Now().Unix()),
			},
			err: errors.New("invalid authority"),
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateAirdrop{
				Authority: authority,
				Intent:    strconv.Itoa(100000),
				Amount:    200,
				Expiry:    uint64(time.Now().Unix()),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			expected := &types.MsgCreateAirdrop{
				Authority: authority,
				Intent:    strconv.Itoa(0),
			}
			_, err := srv.CreateAirdrop(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateAirdrop(ctx, tc.request)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
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
			request: &types.MsgDeleteAirdrop{
				Authority: authority,
				Intent:    strconv.Itoa(0),
			},
		},
		{
			desc: "InvalidSigner",
			request: &types.MsgDeleteAirdrop{
				Authority: "B",
				Intent:    strconv.Itoa(0),
			},
			err: errors.New("invalid authority"),
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteAirdrop{
				Authority: authority,
				Intent:    strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.TokenomicsKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)

			_, err := srv.CreateAirdrop(ctx, &types.MsgCreateAirdrop{
				Authority: authority,
				Intent:    strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteAirdrop(ctx, tc.request)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
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
