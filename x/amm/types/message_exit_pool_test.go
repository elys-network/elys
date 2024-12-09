package types_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestMsgExitPool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgExitPool
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgExitPool{
				Sender:        "invalid_address",
				ShareAmountIn: math.NewInt(100),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgExitPool{
				Sender:        sample.AccAddress(),
				ShareAmountIn: math.NewInt(100),
			},
		},
		{
			name: "Invalid Minimum Amounts Out",
			msg: types.MsgExitPool{
				Sender:        sample.AccAddress(),
				ShareAmountIn: math.NewInt(100),
				MinAmountsOut: sdk.Coins{sdk.Coin{Denom: "uusdc", Amount: math.NewInt(-100)}},
			},
			err: fmt.Errorf("negative coin amount"),
		},
		{
			name: "ShareAmount is Nil",
			msg: types.MsgExitPool{
				Sender:        sample.AccAddress(),
				ShareAmountIn: math.Int{},
				MinAmountsOut: sdk.Coins{sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			err: types.ErrInvalidShareAmountOut,
		},
		{
			name: "ShareAmount is Negative",
			msg: types.MsgExitPool{
				Sender:        sample.AccAddress(),
				ShareAmountIn: math.NewInt(-100),
				MinAmountsOut: sdk.Coins{sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			err: types.ErrInvalidShareAmountOut,
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
