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
		},
		{
			name: "valid address",
			msg: types.MsgSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            nil,
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(10),
				Recipient:         "",
			},
		},
		{
			name: "Invalid recipient address",
			msg: types.MsgSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            nil,
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(1),
				Recipient:         "cosmos1invalid",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "Invalid TokenOutDenom in route",
			msg: types.MsgSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            []types.SwapAmountInRoute{{TokenOutDenom: "invalid denom"}},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(1),
				Recipient:         sample.AccAddress(),
			},
			err: errors.New("invalid denom"),
		},
		{
			name: "Invalid TokenIn",
			msg: types.MsgSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            []types.SwapAmountInRoute{{TokenOutDenom: "uusdc"}},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(-10)},
				TokenOutMinAmount: math.NewInt(1),
				Recipient:         sample.AccAddress(),
			},
			err: errors.New("negative coin amount"),
		},
		{
			name: "Invalid TokenIn amount",
			msg: types.MsgSwapExactAmountIn{
				Sender:            sample.AccAddress(),
				Routes:            []types.SwapAmountInRoute{{TokenOutDenom: "uusdc"}},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(0)},
				TokenOutMinAmount: math.NewInt(1),
				Recipient:         sample.AccAddress(),
			},
			err: errors.New("token in is zero"),
		},
		{
			name: "Invalid routes with same Input and Output denom",
			msg: types.MsgSwapExactAmountIn{
				Sender: sample.AccAddress(),
				Routes: []types.SwapAmountInRoute{
					{TokenOutDenom: "uusdc"},
					{TokenOutDenom: "uusdc"},
					{TokenOutDenom: "uelys"},
				},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(1),
				Recipient:         sample.AccAddress(),
			},
			err: errors.New("has the same input and output denom as the previous route"),
		},
		{
			name: "Valid multiple routes",
			msg: types.MsgSwapExactAmountIn{
				Sender: sample.AccAddress(),
				Routes: []types.SwapAmountInRoute{
					{TokenOutDenom: "uusdc"},
					{TokenOutDenom: "uelys"},
				},
				TokenIn:           sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenOutMinAmount: math.NewInt(1),
				Recipient:         sample.AccAddress(),
			},
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
