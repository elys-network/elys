package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClosePositions_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgClosePositions
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgClosePositions{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgClosePositions{
				Creator: sample.AccAddress(),
				StopLoss: []PositionRequest{
					{Address: sample.AccAddress(),
						Id: 1},
				},
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
