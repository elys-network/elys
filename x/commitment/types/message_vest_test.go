package types_test

import (
	"testing"

	"github.com/elys-network/elys/v7/x/commitment/types"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
)

func TestMsgVest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgVest
		err  error
	}{
		{
			name: "invalid address",
			msg:  types.MsgVest{"invalid_address", math.ZeroInt(), ptypes.Eden},
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  types.MsgVest{sample.AccAddress(), math.OneInt(), ptypes.Eden},
		},
		{
			name: "invalid denom",
			msg:  types.MsgVest{sample.AccAddress(), math.OneInt(), "%%%%"},
			err:  types.ErrInvalidDenom,
		},
		{
			name: "invalid amount - negative",
			msg:  types.MsgVest{sample.AccAddress(), math.NewInt(-200), ptypes.Eden},
			err:  types.ErrInvalidAmount,
		},
		{
			name: "invalid amount - nil",
			msg:  types.MsgVest{sample.AccAddress(), math.Int{}, ptypes.Eden},
			err:  types.ErrInvalidAmount,
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
