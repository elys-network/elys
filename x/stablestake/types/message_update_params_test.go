package types

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdkmath.NewInt(1000),
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdkmath.NewInt(1000),
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
