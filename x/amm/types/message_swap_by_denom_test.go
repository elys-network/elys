package types

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSwapByDenom_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwapByDenom
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwapByDenom{
				Sender: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: MsgSwapByDenom{
				Sender:   sample.AccAddress(),
				Amount:   sdk.NewInt64Coin(ptypes.ATOM, 10),
				DenomIn:  ptypes.ATOM,
				DenomOut: ptypes.BaseCurrency,
			},
		},
		{
			name: "Invalid Amount",
			msg: MsgSwapByDenom{
				Sender:   sample.AccAddress(),
				Amount:   sdk.Coin{Denom: ptypes.ATOM, Amount: sdkmath.NewInt(-10)},
				DenomIn:  ptypes.ATOM,
				DenomOut: ptypes.BaseCurrency,
			},
			err: fmt.Errorf("negative coin amount"),
		},
		{
			name: "Invalid DenomIn",
			msg: MsgSwapByDenom{
				Sender:   sample.AccAddress(),
				Amount:   sdk.Coin{Denom: ptypes.ATOM, Amount: sdkmath.NewInt(10)},
				DenomIn:  "invalid denom in",
				DenomOut: ptypes.BaseCurrency,
			},
			err: fmt.Errorf("invalid denom"),
		},
		{
			name: "Invalid Denomout",
			msg: MsgSwapByDenom{
				Sender:   sample.AccAddress(),
				Amount:   sdk.Coin{Denom: ptypes.ATOM, Amount: sdkmath.NewInt(10)},
				DenomIn:  ptypes.ATOM,
				DenomOut: "invalid denom out",
			},
			err: fmt.Errorf("invalid denom"),
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
