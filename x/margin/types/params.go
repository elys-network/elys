package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyLeverageMax                              = []byte("LeverageMax")
	KeyInterestRateMax                          = []byte("InterestRateMax")
	KeyInterestRateMin                          = []byte("InterestRateMin")
	KeyInterestRateIncrease                     = []byte("InterestRateIncrease")
	KeyInterestRateDecrease                     = []byte("InterestRateDecrease")
	KeyHealthGainFactor                         = []byte("HealthGainFactor")
	KeyEpochLength                              = []byte("EpochLength")
	KeyRemovalQueueThreshold                    = []byte("RemovalQueueThreshold")
	KeyMaxOpenPositions                         = []byte("MaxOpenPositions")
	KeyPoolOpenThreshold                        = []byte("PoolOpenThreshold")
	KeyForceCloseFundPercentage                 = []byte("ForceCloseFundPercentage")
	KeyForceCloseFundAddress                    = []byte("ForceCloseFundAddress")
	KeyIncrementalInterestPaymentFundPercentage = []byte("IncrementalInterestPaymentFundPercentage")
	KeyIncrementalInterestPaymentFundAddress    = []byte("IncrementalInterestPaymentFundAddress")
	KeySqModifier                               = []byte("SqModifier")
	KeySafetyFactor                             = []byte("SafetyFactor")
	KeyIncrementalInterestPaymentEnabled        = []byte("IncrementalInterestPaymentEnabled")
	KeyWhitelistingEnabled                      = []byte("WhitelistingEnabled")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		LeverageMax:                              sdk.NewDec(10),
		InterestRateMax:                          sdk.NewDec(1),
		InterestRateMin:                          sdk.NewDec(1),
		InterestRateIncrease:                     sdk.NewDec(1),
		InterestRateDecrease:                     sdk.NewDec(1),
		HealthGainFactor:                         sdk.NewDec(1),
		EpochLength:                              (int64)(1),
		RemovalQueueThreshold:                    sdk.NewDec(1),
		MaxOpenPositions:                         (int64)(9999),
		PoolOpenThreshold:                        sdk.NewDec(1),
		ForceCloseFundPercentage:                 sdk.NewDec(1),
		ForceCloseFundAddress:                    "",
		IncrementalInterestPaymentFundPercentage: sdk.NewDec(1),
		IncrementalInterestPaymentFundAddress:    "",
		SqModifier:                               sdk.NewDec(1),
		SafetyFactor:                             sdk.NewDec(1),
		IncrementalInterestPaymentEnabled:        true,
		WhitelistingEnabled:                      true,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLeverageMax, &p.LeverageMax, validateLeverageMax),
		paramtypes.NewParamSetPair(KeyInterestRateMax, &p.InterestRateMax, validateInterestRateMax),
		paramtypes.NewParamSetPair(KeyInterestRateMin, &p.InterestRateMin, validateInterestRateMin),
		paramtypes.NewParamSetPair(KeyInterestRateIncrease, &p.InterestRateIncrease, validateInterestRateIncrease),
		paramtypes.NewParamSetPair(KeyInterestRateDecrease, &p.InterestRateDecrease, validateInterestRateDecrease),
		paramtypes.NewParamSetPair(KeyHealthGainFactor, &p.HealthGainFactor, validateHealthGainFactor),
		paramtypes.NewParamSetPair(KeyEpochLength, &p.EpochLength, validateEpochLength),
		paramtypes.NewParamSetPair(KeyRemovalQueueThreshold, &p.RemovalQueueThreshold, validateRemovalQueueThreshold),
		paramtypes.NewParamSetPair(KeyMaxOpenPositions, &p.MaxOpenPositions, validateMaxOpenPositions),
		paramtypes.NewParamSetPair(KeyPoolOpenThreshold, &p.PoolOpenThreshold, validatePoolOpenThreshold),
		paramtypes.NewParamSetPair(KeyForceCloseFundPercentage, &p.ForceCloseFundPercentage, validateForceCloseFundPercentage),
		paramtypes.NewParamSetPair(KeyForceCloseFundAddress, &p.ForceCloseFundAddress, validateForceCloseFundAddress),
		paramtypes.NewParamSetPair(KeyIncrementalInterestPaymentFundPercentage, &p.IncrementalInterestPaymentFundPercentage, validateIncrementalInterestPaymentFundPercentage),
		paramtypes.NewParamSetPair(KeyIncrementalInterestPaymentFundAddress, &p.IncrementalInterestPaymentFundAddress, validateIncrementalInterestPaymentFundAddress),
		paramtypes.NewParamSetPair(KeySqModifier, &p.SqModifier, validateSqModifier),
		paramtypes.NewParamSetPair(KeySafetyFactor, &p.SafetyFactor, validateSafetyFactor),
		paramtypes.NewParamSetPair(KeyIncrementalInterestPaymentEnabled, &p.IncrementalInterestPaymentEnabled, validateIncrementalInterestPaymentEnabled),
		paramtypes.NewParamSetPair(KeyWhitelistingEnabled, &p.WhitelistingEnabled, validateWhitelistingEnabled),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateLeverageMax(p.LeverageMax); err != nil {
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
	if err := validateRemovalQueueThreshold(p.RemovalQueueThreshold); err != nil {
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
	if err := validateIncrementalInterestPaymentFundPercentage(p.IncrementalInterestPaymentFundPercentage); err != nil {
		return err
	}
	if err := validateIncrementalInterestPaymentFundAddress(p.IncrementalInterestPaymentFundAddress); err != nil {
		return err
	}
	if err := validateSqModifier(p.SqModifier); err != nil {
		return err
	}
	if err := validateSafetyFactor(p.SafetyFactor); err != nil {
		return err
	}
	if err := validateIncrementalInterestPaymentEnabled(p.IncrementalInterestPaymentEnabled); err != nil {
		return err
	}
	if err := validateWhitelistingEnabled(p.WhitelistingEnabled); err != nil {
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
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("leverage max must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("leverage max must be positive: %s", v)
	}
	if v.GT(sdk.NewDec(10)) {
		return fmt.Errorf("leverage max too large: %s", v)
	}

	return nil
}

func validateInterestRateMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
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
	v, ok := i.(sdk.Dec)
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
	v, ok := i.(sdk.Dec)
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
	v, ok := i.(sdk.Dec)
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
	v, ok := i.(sdk.Dec)
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
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateRemovalQueueThreshold(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("removal queue threashold must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("removal queue threashold must be positive: %s", v)
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
	v, ok := i.(sdk.Dec)
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

func validateIncrementalInterestPaymentFundPercentage(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("incremental interest payment fund percentage must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("incremental interest payment fund percentage must be positive: %s", v)
	}

	return nil
}

func validateIncrementalInterestPaymentFundAddress(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateSqModifier(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("sq modifier must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("sq modifier must be positive: %s", v)
	}

	return nil
}

func validateSafetyFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
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

func validateIncrementalInterestPaymentEnabled(i interface{}) error {
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
	v, ok := i.(sdk.Dec)
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
