package types_test

import (
	"testing"

	"github.com/elys-network/elys/v4/x/commitment/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v4/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgStake_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgStake
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgStake{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgStake{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(200),
				Asset:   ptypes.ATOM,
			},
		},
		{
			name: "nil amount",
			msg: types.MsgStake{
				Creator: sample.AccAddress(),
				Amount:  math.Int{},
				Asset:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "negative amount",
			msg: types.MsgStake{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-1),
				Asset:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "invalid asset",
			msg: types.MsgStake{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(1),
				Asset:   "",
			},
			err: types.ErrInvalidDenom,
		},
		{
			name: "invalid validator address",
			msg: types.MsgStake{
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
