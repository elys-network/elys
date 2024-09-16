package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/tradeshield/types"
)

func TestPendingPerpetualOrderMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	creator := "A"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreatePendingPerpetualOrder(ctx, &types.MsgCreatePendingPerpetualOrder{Creator: creator})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}

func TestPendingPerpetualOrderMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdatePendingPerpetualOrder
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdatePendingPerpetualOrder{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdatePendingPerpetualOrder{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdatePendingPerpetualOrder{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx := setupMsgServer(t)
			_, err := srv.CreatePendingPerpetualOrder(ctx, &types.MsgCreatePendingPerpetualOrder{Creator: creator})
			require.NoError(t, err)

			_, err = srv.UpdatePendingPerpetualOrder(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPendingPerpetualOrderMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeletePendingPerpetualOrder
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeletePendingPerpetualOrder{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeletePendingPerpetualOrder{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeletePendingPerpetualOrder{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx := setupMsgServer(t)

			_, err := srv.CreatePendingPerpetualOrder(ctx, &types.MsgCreatePendingPerpetualOrder{Creator: creator})
			require.NoError(t, err)
			_, err = srv.DeletePendingPerpetualOrder(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
