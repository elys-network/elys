package types

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgVest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgVest
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgVest{
				Creator: "invalid_address",
				Amount:  sdkmath.ZeroInt(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgVest{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.ZeroInt(),
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
