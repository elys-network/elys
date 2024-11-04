package types_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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
				PoolParams: &types.PoolParams{
					SwapFee:                     sdkmath.LegacyZeroDec(),
					ExitFee:                     sdkmath.LegacyZeroDec(),
					UseOracle:                   false,
					WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
					WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
					WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
					ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
					FeeDenom:                    ptypes.BaseCurrency,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdatePoolParams{
				Authority: sample.AccAddress(),
				PoolParams: &types.PoolParams{
					SwapFee:                     sdkmath.LegacyZeroDec(),
					ExitFee:                     sdkmath.LegacyZeroDec(),
					UseOracle:                   false,
					WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
					WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
					WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
					ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
					FeeDenom:                    ptypes.BaseCurrency,
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
