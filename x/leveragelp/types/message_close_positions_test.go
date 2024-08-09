package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
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
			},
		},
		{
			name: "invalid address liquidations",
			msg: MsgClosePositions{
				Creator: sample.AccAddress(),
				Liquidate: []*PositionRequest{
					{
						Address: "invalid address",
						Id:      uint64(1),
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address liquidations",
			msg: MsgClosePositions{
				Creator: sample.AccAddress(),
				Liquidate: []*PositionRequest{
					{
						Address: sample.AccAddress(),
						Id:      uint64(1),
					},
				},
			},
		},
		{
			name: "invalid address stoploss",
			msg: MsgClosePositions{
				Creator: sample.AccAddress(),
				Stoploss: []*PositionRequest{
					{
						Address: "invalid address",
						Id:      uint64(1),
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address stoploss",
			msg: MsgClosePositions{
				Creator: sample.AccAddress(),
				Stoploss: []*PositionRequest{
					{
						Address: sample.AccAddress(),
						Id:      uint64(1),
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
