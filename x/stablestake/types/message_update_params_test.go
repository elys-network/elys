package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v5/testutil/sample"
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
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
		},
		{
			name: "invalid deposit denom",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative redemption rate",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(-20, 1), // -2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative epoch length",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                -1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(-3, 2), // -0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate min",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(-1, 2), // -0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate increase",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(-2, 2), // -0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative interest rate decrease",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(-1, 2), // -0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1),  // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative health gain factor",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2),  // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2),  // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2),  // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2),  // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(-1, 1), // -0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative total value",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(-1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.7"),
				},
			},
			err: ErrInvalidParams,
		},
		{
			name: "negative max leverage ratio",
			msg: MsgUpdateParams{
				Authority: sample.AccAddress(),
				Params: &Params{
					LegacyDepositDenom:         "stake",
					LegacyRedemptionRate:       sdkmath.LegacyNewDecWithPrec(20, 1), // 2.0
					EpochLength:                1000,
					LegacyInterestRate:         sdkmath.LegacyNewDecWithPrec(3, 2), // 0.03
					LegacyInterestRateMax:      sdkmath.LegacyNewDecWithPrec(5, 2), // 0.05
					LegacyInterestRateMin:      sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyInterestRateIncrease: sdkmath.LegacyNewDecWithPrec(2, 2), // 0.02
					LegacyInterestRateDecrease: sdkmath.LegacyNewDecWithPrec(1, 2), // 0.01
					LegacyHealthGainFactor:     sdkmath.LegacyNewDecWithPrec(1, 1), // 0.1
					TotalValue:                 sdkmath.NewInt(1000),
					LegacyMaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("-0.7"),
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
