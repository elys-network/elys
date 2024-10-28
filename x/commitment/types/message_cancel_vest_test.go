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
			msg: MsgCancelVest{
				Creator: "invalid_address",
				Amount:  math.ZeroInt(),
				Denom:   ptypes.Eden,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgCancelVest{
				Creator: sample.AccAddress(),
				Amount:  math.ZeroInt(),
				Denom:   ptypes.Eden,
			},
		},
		{
			name: "invalid denom",
			msg: MsgCancelVest{
				Creator: sample.AccAddress(),
				Amount:  math.ZeroInt(),
				Denom:   "uusdc",
			},
			err: ErrInvalidDenom,
		},
		{
			name: "invalid amount - nil",
			msg: MsgCancelVest{
				Creator: sample.AccAddress(),
				Denom:   ptypes.Eden,
			},
			err: ErrInvalidAmount,
		},
		{
			name: "invalid amount - negative",
			msg: MsgCancelVest{
				Creator: sample.AccAddress(),
				Denom:   ptypes.Eden,
				Amount:  math.NewInt(-200),
			},
			err: ErrInvalidAmount,
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
