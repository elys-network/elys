package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePendingSpotOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePendingSpotOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePendingSpotOrder{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePendingSpotOrder{
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

func TestMsgUpdatePendingSpotOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdatePendingSpotOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdatePendingSpotOrder{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdatePendingSpotOrder{
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

func TestMsgDeletePendingSpotOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeletePendingSpotOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeletePendingSpotOrder{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeletePendingSpotOrder{
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
