package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateEntry_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateEntry
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateEntry{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateEntry{
				Authority: sample.AccAddress(),
				Decimals:  6,
				BaseDenom: "uusdc",
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

func TestMsgDeleteEntry_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteEntry
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteEntry{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteEntry{
				Authority: sample.AccAddress(),
				BaseDenom: "uusdc",
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
