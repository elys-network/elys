package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/x/tradeshield/types"
)

func TestPendingSpotOrderMsgServerCreate(t *testing.T) {
	srv, ctx := setupMsgServer(t)
	creator := "A"
	for i := 0; i < 5; i++ {
		resp, err := srv.CreatePendingSpotOrder(ctx, &types.MsgCreatePendingSpotOrder{OwnerAddress: creator})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))
	}
}

func TestPendingSpotOrderMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdatePendingSpotOrder
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdatePendingSpotOrder{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdatePendingSpotOrder{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdatePendingSpotOrder{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx := setupMsgServer(t)
			_, err := srv.CreatePendingSpotOrder(ctx, &types.MsgCreatePendingSpotOrder{OwnerAddress: creator})
			require.NoError(t, err)

			_, err = srv.UpdatePendingSpotOrder(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPendingSpotOrderMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeletePendingSpotOrder
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeletePendingSpotOrder{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeletePendingSpotOrder{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeletePendingSpotOrder{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			srv, ctx := setupMsgServer(t)

			_, err := srv.CreatePendingSpotOrder(ctx, &types.MsgCreatePendingSpotOrder{OwnerAddress: creator})
			require.NoError(t, err)
			_, err = srv.DeletePendingSpotOrder(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
