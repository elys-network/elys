package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateGenesisInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateGenesisInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateGenesisInflation{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateGenesisInflation{
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

func TestMsgUpdateGenesisInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateGenesisInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateGenesisInflation{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateGenesisInflation{
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

func TestMsgDeleteGenesisInflation_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteGenesisInflation
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteGenesisInflation{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteGenesisInflation{
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
