package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		DepositDenom:         "uusdc",
		RedemptionRate:       math.LegacyOneDec(),
		EpochLength:          1,
		InterestRate:         math.LegacyMustNewDecFromStr("0.15"),
		InterestRateMax:      math.LegacyMustNewDecFromStr("0.17"),
		InterestRateMin:      math.LegacyMustNewDecFromStr("0.12"),
		InterestRateIncrease: math.LegacyMustNewDecFromStr("0.01"),
		InterestRateDecrease: math.LegacyMustNewDecFromStr("0.01"),
		HealthGainFactor:     math.LegacyOneDec(),
		TotalValue:           math.ZeroInt(),
		MaxLeverageRatio:     math.LegacyMustNewDecFromStr("0.7"),
		MaxWithdrawRatio:     math.LegacyMustNewDecFromStr("0.7"),
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
	if err := sdk.ValidateDenom(p.DepositDenom); err != nil {
		return err
	}

	if err := CheckLegacyDecNilAndNegative(p.RedemptionRate, "RedemptionRate"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.InterestRate, "InterestRate"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.InterestRateMax, "InterestRateMax"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.InterestRateMin, "InterestRateMin"); err != nil {
		return err
	}
	if p.InterestRateMax.LT(p.InterestRateMin) {
		return fmt.Errorf("InterestRateMax must be greater than InterestRateMin")
	}
	if err := CheckLegacyDecNilAndNegative(p.InterestRateIncrease, "InterestRateIncrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.InterestRateDecrease, "InterestRateDecrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.HealthGainFactor, "HealthGainFactor"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.MaxLeverageRatio, "MaxLeverageRatio"); err != nil {
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
