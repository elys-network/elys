package types_test

import (
	"errors"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpFromSwapExactAmountIn_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpFrontSwapExactAmountIn
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpFrontSwapExactAmountIn{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgUpFrontSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            nil,
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(10),
			},
		},
		{
			name: "Invalid TokenOutDenom in route",
			msg: types.MsgUpFrontSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            []types.SwapAmountInRoute{{TokenOutDenom: "invalid denom"}},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(1),
			},
			err: errors.New("invalid denom"),
		},
		{
			name: "Invalid TokenIn",
			msg: types.MsgUpFrontSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            []types.SwapAmountInRoute{{TokenOutDenom: "uusdc"}},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(-10)},
				TokenOutMinAmount: math.NewInt(1),
			},
			err: errors.New("negative coin amount"),
		},
		{
			name: "Invalid TokenIn amount",
			msg: types.MsgUpFrontSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            []types.SwapAmountInRoute{{TokenOutDenom: "uusdc"}},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(0)},
				TokenOutMinAmount: math.NewInt(1),
			},
			err: errors.New("token in is zero"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
