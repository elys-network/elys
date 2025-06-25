package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgTogglePoolEdenRewards_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTogglePoolEdenRewards
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgTogglePoolEdenRewards{
				Authority: "invalid_address",
				PoolId:    1,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgTogglePoolEdenRewards{
				Authority: sample.AccAddress(),
				PoolId:    1,
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
