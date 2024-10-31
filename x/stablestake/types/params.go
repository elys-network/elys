package types

import (
	"fmt"

	"cosmossdk.io/math"
	"gopkg.in/yaml.v2"
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

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateDepositDenom(p.DepositDenom); err != nil {
		return err
	}

	if err := validateRedemptionRate(p.RedemptionRate); err != nil {
		return err
	}
	if err := validateInterestRate(p.InterestRate); err != nil {
		return err
	}
	if err := validateInterestRateMax(p.InterestRateMax); err != nil {
		return err
	}
	if err := validateInterestRateMin(p.InterestRateMin); err != nil {
		return err
	}
	if err := validateInterestRateIncrease(p.InterestRateIncrease); err != nil {
		return err
	}
	if err := validateInterestRateDecrease(p.InterestRateDecrease); err != nil {
		return err
	}
	if err := validateHealthGainFactor(p.HealthGainFactor); err != nil {
		return err
	}
	if err := validateEpochLength(p.EpochLength); err != nil {
		return err
	}
	if err := validateTotalValue(p.TotalValue); err != nil {
		return err
	}
	if err := validateMaxLeverageRatio(p.MaxLeverageRatio); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p LegacyParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateDepositDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("deposit denom should not be empty")
	}

	return nil
}

func validateRedemptionRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid redemption rate type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("redemption rate must be not nil")
	}
	if v.LT(math.LegacyOneDec()) {
		return fmt.Errorf("redemption rate must be bigger than 1: %s", v)
	}

	return nil
}

func validateEpochLength(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("epoch length should be positive: %d", v)
	}

	return nil
}

func validateInterestRate(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest must be positive: %s", v)
	}

	return nil
}

func validateInterestRateMax(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest max must be positive: %s", v)
	}

	return nil
}

func validateInterestRateMin(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest min must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest min must be positive: %s", v)
	}

	return nil
}

func validateInterestRateIncrease(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest rate increase must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest rate increase must be positive: %s", v)
	}

	return nil
}

func validateInterestRateDecrease(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("interest rate decrease must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("interest rate decrease must be positive: %s", v)
	}

	return nil
}

func validateHealthGainFactor(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("health gain factor must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("health gain factor must be positive: %s", v)
	}

	return nil
}

func validateTotalValue(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("total value must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("total value must be positive: %s", v)
	}

	return nil
}

func validateMaxLeverageRatio(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("max leverage percent must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("max leverage percent must be positive: %s", v)
	}

	return nil
}
