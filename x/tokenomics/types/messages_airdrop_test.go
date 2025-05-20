package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateAirdrop_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateAirdrop
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateAirdrop{
				Authority: "invalid_address",
				Intent:    "Airdrop for early adopters",
				Amount:    1000,
				Expiry:    1672531199,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateAirdrop{
				Authority: sample.AccAddress(),
				Intent:    "Airdrop for early adopters",
				Amount:    1000,
				Expiry:    1672531199,
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

func TestMsgUpdateAirdrop_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateAirdrop
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateAirdrop{
				Authority: "invalid_address",
				Intent:    "Airdrop for early adopters",
				Amount:    1000,
				Expiry:    1672531199,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateAirdrop{
				Authority: sample.AccAddress(),
				Intent:    "Airdrop for early adopters",
				Amount:    1000,
				Expiry:    1672531199,
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

func TestMsgDeleteAirdrop_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteAirdrop
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteAirdrop{
				Authority: "invalid_address",
				Intent:    "Airdrop for early adopters",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteAirdrop{
				Authority: sample.AccAddress(),
				Intent:    "Airdrop for early adopters",
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
