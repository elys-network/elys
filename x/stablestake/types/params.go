package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"
)

// NewParams creates a new Params instance
func NewParams(
	depositDenom string,
	redemptionRate math.LegacyDec,
	epochLength int64,
	interestRate math.LegacyDec,
	interestRateMax math.LegacyDec,
	interestRateMin math.LegacyDec,
	interestRateIncrease math.LegacyDec,
	interestRateDecrease math.LegacyDec,
	healthGainFactor math.LegacyDec,
	totalValue math.Int,
	MaxLeveragePercent math.LegacyDec,
) Params {
	return Params{
		DepositDenom:         depositDenom,
		RedemptionRate:       redemptionRate,
		EpochLength:          epochLength,
		InterestRate:         interestRate,
		InterestRateMax:      interestRateMax,
		InterestRateMin:      interestRateMin,
		InterestRateIncrease: interestRateIncrease,
		InterestRateDecrease: interestRateDecrease,
		HealthGainFactor:     healthGainFactor,
		TotalValue:           totalValue,
		MaxLeverageRatio:     MaxLeveragePercent,
	}
}

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
