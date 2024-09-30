package types

import (
	"cosmossdk.io/math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgBond_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBond
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgBond{
				Creator: "invalid_address",
				Amount:  math.NewInt(100),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(100),
			},
		},
		{
			name: "negative amount",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-100),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "zero amount",
			msg: MsgBond{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(0),
			},
			err: sdkerrors.ErrInvalidRequest,
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
