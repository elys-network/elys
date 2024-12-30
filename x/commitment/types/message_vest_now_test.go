package types_test

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgVestNow_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgVestNow
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgVestNow{
				Creator: "invalid_address",
				Amount:  sdkmath.OneInt(),
				Denom:   ptypes.ATOM,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgVestNow{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.OneInt(),
				Denom:   ptypes.ATOM,
			},
		},
		{
			name: "amount is nil",
			msg: types.MsgVestNow{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.Int{},
				Denom:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "amount is -ve",
			msg: types.MsgVestNow{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.NewInt(-14),
				Denom:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "invalid denom",
			msg: types.MsgVestNow{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.NewInt(14),
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
