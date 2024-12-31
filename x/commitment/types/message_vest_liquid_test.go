package types_test

import (
	"errors"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgVestLiquid_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgVestLiquid
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgVestLiquid{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgVestLiquid{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(200),
				Denom:   ptypes.ATOM,
			},
		},
		{
			name: "amount is nil",
			msg: types.MsgVestLiquid{
				Creator: sample.AccAddress(),
				Amount:  math.Int{},
				Denom:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "amount is -ve",
			msg: types.MsgVestLiquid{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(-14),
				Denom:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "invalid denom",
			msg: types.MsgVestLiquid{
				Creator: sample.AccAddress(),
				Amount:  math.NewInt(14),
				Denom:   "",
			},
			err: errors.New("invalid denom"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
