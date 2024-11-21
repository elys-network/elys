package types_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSwapExactAmountIn_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgSwapExactAmountIn
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgSwapExactAmountIn{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            nil,
				TokenIn:           sdk.Coin{ptypes.ATOM, math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(10),
				Recipient:         "",
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
