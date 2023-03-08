package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateTimeBasedInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateTimeBasedInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateTimeBasedInflation{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateTimeBasedInflation{
				Authority: sample.AccAddress(),
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

func TestMsgUpdateTimeBasedInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateTimeBasedInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateTimeBasedInflation{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateTimeBasedInflation{
				Authority: sample.AccAddress(),
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

func TestMsgDeleteTimeBasedInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteTimeBasedInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteTimeBasedInflation{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteTimeBasedInflation{
				Authority: sample.AccAddress(),
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
