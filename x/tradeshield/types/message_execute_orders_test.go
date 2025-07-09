package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgExecuteOrders_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgExecuteOrders
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgExecuteOrders{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgExecuteOrders{
				Creator:      sample.AccAddress(),
				SpotOrderIds: []uint64{1},
				PerpetualOrders: []PerpetualOrderKey{
					{
						OwnerAddress: sample.AccAddress(),
						PoolId:       1,
						OrderId:      1,
					},
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
