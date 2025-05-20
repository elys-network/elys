package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSetPortfolio_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetPortfolio
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetPortfolio{
				Creator: "invalid_address",
				User:    sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetPortfolio{
				Creator: sample.AccAddress(),
				User:    sample.AccAddress(),
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
