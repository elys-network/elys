package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSetWithdrawAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetWithdrawAddress
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetWithdrawAddress{
				DelegatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetWithdrawAddress{
				DelegatorAddress: sample.AccAddress(),
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

func TestMsgWithdrawValidatorCommission_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgWithdrawValidatorCommission
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgWithdrawValidatorCommission{
				ValidatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgWithdrawValidatorCommission{
				ValidatorAddress: sample.AccAddress(),
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

func TestMsgWithdrawDelegatorReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgWithdrawDelegatorReward
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgWithdrawDelegatorReward{
				DelegatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgWithdrawDelegatorReward{
				DelegatorAddress: sample.AccAddress(),
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
