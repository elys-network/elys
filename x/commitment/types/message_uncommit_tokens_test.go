package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUncommitTokens_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUncommitTokens
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUncommitTokens{
				Creator: "invalid_address",
				Amount:  sdkmath.ZeroInt(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUncommitTokens{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.OneInt(),
				Denom:   ptypes.ATOM,
			},
		},
		{
			name: "nil amount",
			msg: types.MsgUncommitTokens{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.Int{},
				Denom:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "negative amount",
			msg: types.MsgUncommitTokens{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.NewInt(-1),
				Denom:   ptypes.ATOM,
			},
			err: types.ErrInvalidAmount,
		},
		{
			name: "invalid denom",
			msg: types.MsgUncommitTokens{
				Creator: sample.AccAddress(),
				Amount:  sdkmath.NewInt(1),
				Denom:   "",
			},
			err: types.ErrInvalidDenom,
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
