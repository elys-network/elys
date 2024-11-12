package types

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/stretchr/testify/assert"
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
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
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
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
		},
		{
			name: "invalid deposit denom",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative redemption rate",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(-20, 1), // -2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative epoch length",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          -1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(-3, 2), // -0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate min",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(-1, 2), // -0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate increase",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(-2, 2), // -0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate decrease",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(-1, 2), // -0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative health gain factor",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					DepositDenom:         "stake",
					RedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:          1000,
					InterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					InterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					InterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					InterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					InterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					HealthGainFactor:     sdkmath.LegacyNewDecWithPrec(-1, 1), // -0.1
					TotalValue:           sdkmath.NewInt(1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative total value",
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
					TotalValue:           sdkmath.NewInt(-1000),
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative max leverage ratio",
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
					MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("-0.7"),
				},
			},
			err: ErrInvalidParams,
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

func TestNewMsgUpdateParams(t *testing.T) {

	accAdress := sample.AccAddress()

	params := DefaultParams()
	got := NewMsgUpdateParams(
		accAdress,
		&params,
	)

	want := &MsgUpdateParams{
		Authority: accAdress,
		Params:    &params,
	}

	assert.Equal(t, want, got)
}
