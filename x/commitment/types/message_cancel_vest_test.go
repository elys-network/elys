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
			msg:  *NewMsgCancelVest("invalid_address", math.ZeroInt(), ptypes.Eden),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  *NewMsgCancelVest(sample.AccAddress(), math.ZeroInt(), ptypes.Eden),
		},
		{
			name: "invalid denom",
			msg:  *NewMsgCancelVest(sample.AccAddress(), math.ZeroInt(), "invalid_denom"),
			err:  ErrInvalidDenom,
		},
		{
			name: "invalid amount - negative",
			msg:  *NewMsgCancelVest(sample.AccAddress(), math.NewInt(-200), ptypes.Eden),
			err:  ErrInvalidAmount,
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
