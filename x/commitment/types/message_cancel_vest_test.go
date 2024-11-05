package types

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCancelVest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCancelVest
		err  error
	}{
		{
			name: "invalid address",
			msg:  *NewMsgCancelVest("invalid_address", sdk.ZeroInt(), ptypes.Eden),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  *NewMsgCancelVest(sample.AccAddress(), sdk.ZeroInt(), ptypes.Eden),
		},
		{
			name: "invalid denom",
			msg:  *NewMsgCancelVest(sample.AccAddress(), sdk.ZeroInt(), "invalid_denom"),
			err:  ErrInvalidDenom,
		},
		{
			name: "invalid amount - negative",
			msg:  *NewMsgCancelVest(sample.AccAddress(), sdk.NewInt(-200), ptypes.Eden),
			err:  ErrInvalidAmount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.msg.Route(), RouterKey)
			require.Equal(t, tt.msg.Type(), TypeMsgCancelVest)
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			} else {
				creator := sdk.MustAccAddressFromBech32(tt.msg.Creator)
				require.Equal(t, tt.msg.GetSigners(), []sdk.AccAddress{creator})
			}
			require.NoError(t, err)
		})
	}
}

// GetSigners panics
func TestMsgCancelVest_GetSigners(t *testing.T) {
	msg := NewMsgCancelVest("invalid_address", sdk.ZeroInt(), ptypes.Eden)
	require.Panics(t, func() { msg.GetSigners() })
}

// GetSignBytes
func TestMsgCancelVest_GetSignBytes(t *testing.T) {
	msg := NewMsgCancelVest(sample.AccAddress(), sdk.ZeroInt(), ptypes.Eden)
	bz := ModuleCdc.MustMarshalJSON(msg)
	require.Equal(t, msg.GetSignBytes(), sdk.MustSortJSON(bz))
}
