package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/testutil/sample"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddPool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name   string
		msg    types.MsgAddPool
		errMsg string
		err    bool
	}{
		{
			name: "valid message",
			msg: types.MsgAddPool{
				Authority:            sample.AccAddress(),
				DepositDenom:         "ustake",
				InterestRate:         math.LegacyMustNewDecFromStr("0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
			},
			err:    false,
			errMsg: "",
		},
		{
			name: "invalid authority address",
			msg: types.MsgAddPool{
				Authority:            "invalid_address",
				DepositDenom:         "stake",
				InterestRate:         math.LegacyMustNewDecFromStr("0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
			},
			err:    true,
			errMsg: "invalid address",
		},
		{
			name: "invalid deposit denom",
			msg: types.MsgAddPool{
				Authority:            sample.AccAddress(),
				DepositDenom:         "invalid denom",
				InterestRate:         math.LegacyMustNewDecFromStr("0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
			},
			err:    true,
			errMsg: "invalid denom: invalid denom",
		},
		{
			name: "negative interest rate",
			msg: types.MsgAddPool{
				Authority:            sample.AccAddress(),
				DepositDenom:         "stake",
				InterestRate:         math.LegacyMustNewDecFromStr("-0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
			},
			err:    true,
			errMsg: "InterestRate is negative",
		},
		{
			name: "interest rate max less than interest rate min",
			msg: types.MsgAddPool{
				Authority:            sample.AccAddress(),
				DepositDenom:         "stake",
				InterestRate:         math.LegacyMustNewDecFromStr("0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
			},
			err:    true,
			errMsg: "InterestRateMax must be greater than InterestRateMin",
		},
		{
			name: "negative health gain factor",
			msg: types.MsgAddPool{
				Authority:            sample.AccAddress(),
				DepositDenom:         "stake",
				InterestRate:         math.LegacyMustNewDecFromStr("0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("-0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
			},
			err:    true,
			errMsg: "HealthGainFactor is negative",
		},
		{
			name: "negative max leverage ratio",
			msg: types.MsgAddPool{
				Authority:            sample.AccAddress(),
				DepositDenom:         "stake",
				InterestRate:         math.LegacyMustNewDecFromStr("0.03"),
				InterestRateMax:      math.LegacyMustNewDecFromStr("0.05"),
				InterestRateMin:      math.LegacyMustNewDecFromStr("0.01"),
				InterestRateIncrease: math.LegacyMustNewDecFromStr("0.02"),
				InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     math.LegacyMustNewDecFromStr("0.1"),
				MaxLeverageRatio:     math.LegacyMustNewDecFromStr("-0.7"),
			},
			err:    true,
			errMsg: "MaxLeverageRatio is negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err {
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
