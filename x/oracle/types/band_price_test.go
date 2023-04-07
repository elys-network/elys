package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgRequestBandPrice_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRequestBandPrice
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRequestBandPrice{
				Creator:       "invalid_address",
				SourceChannel: "channel-0",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid source channel",
			msg: MsgRequestBandPrice{
				Creator:       sample.AccAddress(),
				SourceChannel: "",
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "valid message",
			msg: MsgRequestBandPrice{
				Creator:       sample.AccAddress(),
				SourceChannel: "channel-0",
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
