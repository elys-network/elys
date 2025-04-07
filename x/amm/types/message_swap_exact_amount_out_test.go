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

func TestMsgSwapExactAmountOut_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgSwapExactAmountOut
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgSwapExactAmountOut{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgSwapExactAmountOut{
				Sender:           sample.AccAddress(),
				Routes:           nil,
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        "",
			},
		},
		{
			name: "Invalid recipient address",
			msg: types.MsgSwapExactAmountOut{
				Sender:           sample.AccAddress(),
				Routes:           nil,
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        "cosmos1invalid",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "Invalid tokenInDenom in route",
			msg: types.MsgSwapExactAmountOut{
				Sender:           sample.AccAddress(),
				Routes:           []types.SwapAmountOutRoute{{TokenInDenom: "invalid denom"}},
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
			},
			err: errors.New("invalid denom"),
		},
		{
			name: "Invalid tokenOut",
			msg: types.MsgSwapExactAmountOut{
				Sender:           sample.AccAddress(),
				Routes:           []types.SwapAmountOutRoute{{TokenInDenom: "uusdc"}},
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(-10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
			},
			err: errors.New("negative coin amount"),
		},
		{
			name: "Invalid tokenOut amount",
			msg: types.MsgSwapExactAmountOut{
				Sender:           sample.AccAddress(),
				Routes:           []types.SwapAmountOutRoute{{TokenInDenom: "uusdc"}},
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(0)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
			},
			err: errors.New("token in is zero"),
		},
		{
			name: "Invalid routes with same Input and Output denom",
			msg: types.MsgSwapExactAmountOut{
				Sender: sample.AccAddress(),
				Routes: []types.SwapAmountOutRoute{
					{TokenInDenom: "uusdc"},
					{TokenInDenom: "uusdc"},
					{TokenInDenom: "uelys"},
				},
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
			},
			err: errors.New("has the same input and output denom as the previous route"),
		},
		{
			name: "Duplicate TokenInDenom in routes",
			msg: types.MsgSwapExactAmountOut{
				Sender: sample.AccAddress(),
				Routes: []types.SwapAmountOutRoute{
					{TokenInDenom: "uusdc"},
					{TokenInDenom: "uosmo"},
					{TokenInDenom: "uelys"},
					{TokenInDenom: "uusdc"},
				},
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
			},
			err: errors.New("duplicate TokenInDenom found in route 3"),
		},
		{
			name: "Circular swap detected",
			msg: types.MsgSwapExactAmountOut{
				Sender: sample.AccAddress(),
				Routes: []types.SwapAmountOutRoute{
					{TokenInDenom: "uusdc"},
					{TokenInDenom: "uelys"},
				},
				TokenOut:         sdk.Coin{Denom: "uusdc", Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
			},
			err: errors.New("circular swap detected: token in denom matches the last route's token out denom"),
		},
		{
			name: "Valid multiple routes",
			msg: types.MsgSwapExactAmountOut{
				Sender: sample.AccAddress(),
				Routes: []types.SwapAmountOutRoute{
					{TokenInDenom: "uusdc"},
					{TokenInDenom: "uelys"},
				},
				TokenOut:         sdk.Coin{Denom: ptypes.ATOM, Amount: math.NewInt(10)},
				TokenInMaxAmount: math.NewInt(1),
				Recipient:        sample.AccAddress(),
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
