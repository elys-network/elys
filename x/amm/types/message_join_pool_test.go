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

func TestMsgJoinPool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgJoinPool
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgJoinPool{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgJoinPool{
				Sender:         sample.AccAddress(),
				ShareAmountOut: math.NewInt(100),
			},
		},
		{
			name: "Invalid Maximum Amounts in",
			msg: types.MsgJoinPool{
				Sender:         sample.AccAddress(),
				ShareAmountOut: math.NewInt(100),
				MaxAmountsIn:   sdk.Coins{sdk.Coin{Denom: "uusdc", Amount: math.NewInt(-100)}},
			},
			err: fmt.Errorf("negative coin amount"),
		},
		{
			name: "ShareAmount is Nil",
			msg: types.MsgJoinPool{
				Sender:         sample.AccAddress(),
				ShareAmountOut: math.Int{},
				MaxAmountsIn:   sdk.Coins{sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)}},
			},
			err: types.ErrInvalidShareAmountOut,
		},
		{
			name: "ShareAmount is Negative",
			msg: types.MsgJoinPool{
				Sender:         sample.AccAddress(),
				ShareAmountOut: math.NewInt(-100),
				MaxAmountsIn:   sdk.Coins{sdk.Coin{Denom: "uusdc", Amount: math.NewInt(100)}},
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
