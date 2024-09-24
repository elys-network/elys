package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePendingPerpetualOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePendingPerpetualOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePendingPerpetualOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePendingPerpetualOrder{
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

func TestMsgUpdatePendingPerpetualOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdatePendingPerpetualOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdatePendingPerpetualOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdatePendingPerpetualOrder{
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

func TestMsgDeletePendingPerpetualOrder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeletePendingPerpetualOrder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeletePendingPerpetualOrder{
				OwnerAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeletePendingPerpetualOrder{
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
