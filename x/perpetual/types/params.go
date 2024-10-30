package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		PerpetualSwapFee:                               math.LegacyMustNewDecFromStr("0.001"), // 0.1%
		FixedFundingRate:                               math.LegacyMustNewDecFromStr("0.5"),   // 50%
		BorrowInterestRateDecrease:                     math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateIncrease:                     math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateMax:                          math.LegacyMustNewDecFromStr("0.3"),
		BorrowInterestRateMin:                          math.LegacyMustNewDecFromStr("0.1"),
		ForceCloseFundAddress:                          authtypes.NewModuleAddress("zero").String(),
		ForceCloseFundPercentage:                       math.LegacyOneDec(),
		HealthGainFactor:                               math.LegacyMustNewDecFromStr("0.000000220000000000"),
		IncrementalBorrowInterestPaymentEnabled:        true,
		IncrementalBorrowInterestPaymentFundAddress:    authtypes.NewModuleAddress("zero").String(),
		IncrementalBorrowInterestPaymentFundPercentage: math.LegacyMustNewDecFromStr("0.1"),
		LeverageMax:                                    math.LegacyNewDec(25),
		MaxOpenPositions:                               (int64)(3000),
		PoolOpenThreshold:                              math.LegacyMustNewDecFromStr("0.65"),
		SafetyFactor:                                   math.LegacyMustNewDecFromStr("1.050000000000000000"), // 5%
		WhitelistingEnabled:                            false,
		MaxLimitOrder:                                  (int64)(500),
		MinimumLongTakeProfitPriceRatio:                math.LegacyMustNewDecFromStr("1.02"),
		MaximumLongTakeProfitPriceRatio:                math.LegacyMustNewDecFromStr("11"),
		MaximumShortTakeProfitPriceRatio:               math.LegacyMustNewDecFromStr("0.98"),
		EnableTakeProfitCustodyLiabilities:             false,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateLeverageMax(p.LeverageMax); err != nil {
		return err
	}
	if err := validateBorrowInterestRateMax(p.BorrowInterestRateMax); err != nil {
		return err
	}
	if err := validateBorrowInterestRateMin(p.BorrowInterestRateMin); err != nil {
		return err
	}
	if err := validateBorrowInterestRateIncrease(p.BorrowInterestRateIncrease); err != nil {
		return err
	}
	if err := validateBorrowInterestRateDecrease(p.BorrowInterestRateDecrease); err != nil {
		return err
	}
	if err := validateHealthGainFactor(p.HealthGainFactor); err != nil {
		return err
	}
	if err := validateMaxOpenPositions(p.MaxOpenPositions); err != nil {
		return err
	}
	if err := validatePoolOpenThreshold(p.PoolOpenThreshold); err != nil {
		return err
	}
	if err := validateForceCloseFundPercentage(p.ForceCloseFundPercentage); err != nil {
		return err
	}
	if err := validateForceCloseFundAddress(p.ForceCloseFundAddress); err != nil {
		return err
	}
	if err := validateIncrementalBorrowInterestPaymentFundPercentage(p.IncrementalBorrowInterestPaymentFundPercentage); err != nil {
		return err
	}
	if err := validateIncrementalBorrowInterestPaymentFundAddress(p.IncrementalBorrowInterestPaymentFundAddress); err != nil {
		return err
	}
	if err := validateSafetyFactor(p.SafetyFactor); err != nil {
		return err
	}
	if err := validateIncrementalBorrowInterestPaymentEnabled(p.IncrementalBorrowInterestPaymentEnabled); err != nil {
		return err
	}
	if err := validateWhitelistingEnabled(p.WhitelistingEnabled); err != nil {
		return err
	}
	if err := validateFixedFundingRate(p.FixedFundingRate); err != nil {
		return err
	}
	if err := validateSwapFee(p.PerpetualSwapFee); err != nil {
		return err
	}
	if err := validateMaxLimitOrder(p.MaxLimitOrder); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateLeverageMax(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("leverage max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("leverage max must be positive: %s", v)
	}
	return nil
}

func validateBorrowInterestRateMax(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest max must be positive: %s", v)
	}

	return nil
}

func validateBorrowInterestRateMin(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest min must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest min must be positive: %s", v)
	}

	return nil
}

func validateBorrowInterestRateIncrease(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest rate increase must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest rate increase must be positive: %s", v)
	}

	return nil
}

func validateBorrowInterestRateDecrease(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("borrow interest rate decrease must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("borrow interest rate decrease must be positive: %s", v)
	}

	return nil
}

func validateHealthGainFactor(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
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

func validateMaxOpenPositions(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateForceCloseFundPercentage(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("force close fund percentage must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("force close fund percentage must be positive: %s", v)
	}

	return nil
}

func validateForceCloseFundAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateIncrementalBorrowInterestPaymentFundPercentage(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("incremental borrow interest payment fund percentage must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("incremental borrow interest payment fund percentage must be positive: %s", v)
	}

	return nil
}

func validateIncrementalBorrowInterestPaymentFundAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateSafetyFactor(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("safety factor must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("safety factor must be positive: %s", v)
	}

	return nil
}

func validateIncrementalBorrowInterestPaymentEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateWhitelistingEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validatePoolOpenThreshold(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("pool open threshold must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("pool open threshold must be positive: %s", v)
	}

	return nil
}

func validateFixedFundingRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("fixed funding fee must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("fixed funding fee must be positive: %s", v)
	}

	return nil
}

func validateSwapFee(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("swap fee must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("swap fee must be positive: %s", v)
	}

	return nil
}

func validateMaxLimitOrder(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("MaxLimitOrder should not be -ve: %d", v)
	}
	return nil
}
