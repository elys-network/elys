package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	epochtypes "github.com/elys-network/elys/x/epochs/types"
	"gopkg.in/yaml.v2"
)

var ZeroAddress = authtypes.NewModuleAddress("zero").String()

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		SwapFee:                                        sdkmath.LegacyNewDecWithPrec(1, 3), // 0.1%
		FundingFeeCollectionAddress:                    ZeroAddress,
		FundingFeeMinRate:                              sdkmath.LegacyNewDecWithPrec(-1, 3), // -0.1%
		FundingFeeMaxRate:                              sdkmath.LegacyNewDecWithPrec(1, 3),  // 0.1%
		FundingFeeBaseRate:                             sdkmath.LegacyNewDecWithPrec(3, 4),  // 0.03%
		TakeProfitBorrowInterestRateMin:                sdkmath.LegacyOneDec(),
		BorrowInterestRateDecrease:                     sdkmath.LegacyNewDecWithPrec(33, 10),
		BorrowInterestRateIncrease:                     sdkmath.LegacyNewDecWithPrec(33, 10),
		BorrowInterestRateMax:                          sdkmath.LegacyNewDecWithPrec(27, 7),
		BorrowInterestRateMin:                          sdkmath.LegacyNewDecWithPrec(3, 8),
		MinBorrowInterestAmount:                        sdkmath.NewInt(5_000_000),
		EpochLength:                                    (int64)(1),
		ForceCloseFundAddress:                          ZeroAddress,
		ForceCloseFundPercentage:                       sdkmath.LegacyOneDec(),
		HealthGainFactor:                               sdkmath.LegacyNewDecWithPrec(22, 8),
		IncrementalBorrowInterestPaymentEnabled:        true,
		IncrementalBorrowInterestPaymentFundAddress:    ZeroAddress,
		IncrementalBorrowInterestPaymentFundPercentage: sdkmath.LegacyNewDecWithPrec(35, 1), // 35%
		InvariantCheckEpoch:                            epochtypes.DayEpochID,
		LeverageMax:                                    sdkmath.LegacyNewDec(10),
		MaxOpenPositions:                               (int64)(9999),
		PoolOpenThreshold:                              sdkmath.LegacyOneDec(),
		SafetyFactor:                                   sdkmath.LegacyMustNewDecFromStr("1.050000000000000000"), // 5%
		WhitelistingEnabled:                            false,
		MaxLimitOrder:                                  (int64)(500),
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
	if err := validateEpochLength(p.EpochLength); err != nil {
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
	if err := validateInvariantCheckEpoch(p.InvariantCheckEpoch); err != nil {
		return err
	}
	if err := validateTakeProfitBorrowInterestRateMin(p.TakeProfitBorrowInterestRateMin); err != nil {
		return err
	}
	if err := validateFundingFeeBaseRate(p.FundingFeeBaseRate); err != nil {
		return err
	}
	if err := validateFundingFeeMinRate(p.FundingFeeMinRate); err != nil {
		return err
	}
	if err := validateFundingFeeMaxRate(p.FundingFeeMaxRate); err != nil {
		return err
	}
	if err := validateFundingFeeCollectionAddress(p.FundingFeeCollectionAddress); err != nil {
		return err
	}
	if err := validateSwapFee(p.SwapFee); err != nil {
		return err
	}
	if err := validateMinBorrowInterestAmount(p.MinBorrowInterestAmount); err != nil {
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

func (p LegacyParams) String() string {
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

func validateInvariantCheckEpoch(i interface{}) error {
	epoch, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if epoch != epochtypes.DayEpochID && epoch != epochtypes.WeekEpochID && epoch != epochtypes.HourEpochID {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateTakeProfitBorrowInterestRateMin(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("take profit borrow interest rate min must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("take profit borrow interest rate min must be positive: %s", v)
	}

	return nil
}

func validateFundingFeeBaseRate(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("funding fee base rate must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("funding fee base rate must be positive: %s", v)
	}

	return nil
}

func validateFundingFeeMinRate(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("funding fee min rate must be not nil")
	}
	if v.IsPositive() {
		return fmt.Errorf("funding fee min rate must be negative: %s", v)
	}

	return nil
}

func validateFundingFeeMaxRate(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("funding fee max rate must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("funding fee max rate must be positive: %s", v)
	}

	return nil
}

func validateFundingFeeCollectionAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
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

func validateMinBorrowInterestAmount(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("MinBorrowInterestAmount must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("MinBorrowInterestAmount must be positive: %s", v)
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
