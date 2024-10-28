package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateSpotOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateSpotOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateSpotOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateSpotOrder{
				OwnerAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateSpotOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateSpotOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateSpotOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateSpotOrder{
				OwnerAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgCancelSpotOrders_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCancelSpotOrders
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCancelSpotOrders{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCancelSpotOrders{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
