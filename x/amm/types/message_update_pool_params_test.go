package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/testutil/sample"
	"github.com/elys-network/elys/v7/x/amm/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdatePoolParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdatePoolParams
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdatePoolParams{
				Authority: "invalid_address",
				PoolParams: types.PoolParams{
					SwapFee:   math.LegacyZeroDec(),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdatePoolParams{
				Authority: sample.AccAddress(),
				PoolId:    1,
				PoolParams: types.PoolParams{
					SwapFee:   math.LegacyZeroDec(),
					UseOracle: false,
					FeeDenom:  ptypes.BaseCurrency,
				},
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
