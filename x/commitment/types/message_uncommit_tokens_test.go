package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUncommitTokens_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUncommitTokens
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUncommitTokens{
				Creator: "invalid_address",
				Amount:  sdk.ZeroInt(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUncommitTokens{
				Creator: sample.AccAddress(),
				Amount:  sdk.ZeroInt(),
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
