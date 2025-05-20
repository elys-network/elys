package types_test

import (
	"testing"

	"github.com/elys-network/elys/v4/x/commitment/types"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v4/testutil/sample"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCancelVest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgCancelVest
		err  error
	}{
		{
			name: "invalid address",
			msg:  *types.NewMsgCancelVest("invalid_address", math.ZeroInt(), ptypes.Eden),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg:  *types.NewMsgCancelVest(sample.AccAddress(), math.OneInt(), ptypes.Eden),
		},
		{
			name: "invalid denom",
			msg:  *types.NewMsgCancelVest(sample.AccAddress(), math.OneInt(), "invalid_denom"),
			err:  types.ErrInvalidDenom,
		},
		{
			name: "invalid amount - negative",
			msg:  *types.NewMsgCancelVest(sample.AccAddress(), math.NewInt(-200), ptypes.Eden),
			err:  types.ErrInvalidAmount,
		},
		{
			name: "invalid amount - nil",
			msg:  types.MsgCancelVest{sample.AccAddress(), math.Int{}, ptypes.Eden},
			err:  types.ErrInvalidAmount,
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
