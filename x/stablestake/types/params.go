package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		LegacyDepositDenom:         "uusdc",
		LegacyRedemptionRate:       math.LegacyOneDec(),
		EpochLength:                1,
		LegacyInterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		LegacyInterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		LegacyInterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		LegacyInterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		LegacyInterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		LegacyHealthGainFactor:     math.LegacyOneDec(),
		TotalValue:                 math.ZeroInt(),
		LegacyMaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		LegacyMaxWithdrawRatio:     math.LegacyMustNewDecFromStr("0.7"),
	}
}

func CheckLegacyDecNilAndNegative(value math.LegacyDec, name string) error {
	if value.IsNil() {
		return fmt.Errorf("%s is nil", name)
	}
	if value.IsNegative() {
		return fmt.Errorf("%s is negative", name)
	}
	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := sdk.ValidateDenom(p.LegacyDepositDenom); err != nil {
		return err
	}

	if err := CheckLegacyDecNilAndNegative(p.LegacyRedemptionRate, "RedemptionRate"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyInterestRate, "InterestRate"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyInterestRateMax, "InterestRateMax"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyInterestRateMin, "InterestRateMin"); err != nil {
		return err
	}
	if p.LegacyInterestRateMax.LT(p.LegacyInterestRateMin) {
		return fmt.Errorf("InterestRateMax must be greater than InterestRateMin")
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyInterestRateIncrease, "InterestRateIncrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyInterestRateDecrease, "InterestRateDecrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyHealthGainFactor, "HealthGainFactor"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LegacyMaxLeverageRatio, "MaxLeverageRatio"); err != nil {
		return err
	}

	if p.TotalValue.IsNil() {
		return fmt.Errorf("TotalValue is nil")
	}
	if p.TotalValue.IsNegative() {
		return fmt.Errorf("TotalValue is negative")
	}

	if p.EpochLength < 0 {
		return fmt.Errorf("EpochLength is negative")
	}
	return nil
}
