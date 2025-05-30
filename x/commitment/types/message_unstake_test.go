package types_test

import (
	"testing"

	"github.com/elys-network/elys/v6/x/commitment/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v6/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUnstake_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUnstake
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUnstake{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUnstake{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(200),
				Asset:   ptypes.ATOM,
			},
		},
		{
			name: "nil amount",
			msg: types.MsgUnstake{
				Creator: sample.AccAddress(),
				Amount:  math.Int{},
				Asset:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "negative amount",
			msg: types.MsgUnstake{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-1),
				Asset:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "invalid asset",
			msg: types.MsgUnstake{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(1),
				Asset:   "",
			},
			err: types.ErrInvalidDenom,
		},
		{
			name: "invalid validator address",
			msg: types.MsgUnstake{
				Creator:          sample.AccAddress(),
				Amount:           math.NewInt(1),
				Asset:            ptypes.Elys,
				ValidatorAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
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
