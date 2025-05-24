package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSetPriceFeeder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetPriceFeeder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetPriceFeeder{
				Feeder: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetPriceFeeder{
				Feeder: sample.AccAddress(),
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

func TestMsgDeletePriceFeeder_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeletePriceFeeder
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeletePriceFeeder{
				Feeder: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeletePriceFeeder{
				Feeder: sample.AccAddress(),
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
