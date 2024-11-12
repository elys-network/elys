package types

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgFeedPrice_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgFeedPrice
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgFeedPrice{
				Provider: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgFeedPrice{
				Provider: sample.AccAddress(),
				Price:    sdkmath.LegacyMustNewDecFromStr("100"),
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
