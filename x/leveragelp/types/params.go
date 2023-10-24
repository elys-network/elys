package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	epochtypes "github.com/elys-network/elys/x/epochs/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyLeverageMax              = []byte("LeverageMax")
	KeyEpochLength              = []byte("EpochLength")
	KeyRemovalQueueThreshold    = []byte("RemovalQueueThreshold")
	KeyMaxOpenPositions         = []byte("MaxOpenPositions")
	KeyPoolOpenThreshold        = []byte("PoolOpenThreshold")
	KeyForceCloseFundPercentage = []byte("ForceCloseFundPercentage")
	KeyForceCloseFundAddress    = []byte("ForceCloseFundAddress")
	KeySqModifier               = []byte("SqModifier")
	KeySafetyFactor             = []byte("SafetyFactor")
	KeyWhitelistingEnabled      = []byte("WhitelistingEnabled")
	KeyInvariantCheckEpoch      = []byte("InvariantCheckEpoch")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		LeverageMax:              sdk.NewDec(10),
		EpochLength:              (int64)(1),
		RemovalQueueThreshold:    sdk.NewDec(1),
		MaxOpenPositions:         (int64)(9999),
		PoolOpenThreshold:        sdk.NewDec(1),
		ForceCloseFundPercentage: sdk.NewDec(1),
		ForceCloseFundAddress:    "",
		SqModifier:               sdk.NewDec(1),
		SafetyFactor:             sdk.NewDec(1),
		WhitelistingEnabled:      false,
		InvariantCheckEpoch:      epochtypes.DayEpochID,
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
		paramtypes.NewParamSetPair(KeyEpochLength, &p.EpochLength, validateEpochLength),
		paramtypes.NewParamSetPair(KeyRemovalQueueThreshold, &p.RemovalQueueThreshold, validateRemovalQueueThreshold),
		paramtypes.NewParamSetPair(KeyMaxOpenPositions, &p.MaxOpenPositions, validateMaxOpenPositions),
		paramtypes.NewParamSetPair(KeyPoolOpenThreshold, &p.PoolOpenThreshold, validatePoolOpenThreshold),
		paramtypes.NewParamSetPair(KeyForceCloseFundPercentage, &p.ForceCloseFundPercentage, validateForceCloseFundPercentage),
		paramtypes.NewParamSetPair(KeyForceCloseFundAddress, &p.ForceCloseFundAddress, validateForceCloseFundAddress),
		paramtypes.NewParamSetPair(KeySqModifier, &p.SqModifier, validateSqModifier),
		paramtypes.NewParamSetPair(KeySafetyFactor, &p.SafetyFactor, validateSafetyFactor),
		paramtypes.NewParamSetPair(KeyWhitelistingEnabled, &p.WhitelistingEnabled, validateWhitelistingEnabled),
		paramtypes.NewParamSetPair(KeyInvariantCheckEpoch, &p.InvariantCheckEpoch, validateInvariantCheckEpoch),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateLeverageMax(p.LeverageMax); err != nil {
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
	if err := validateSqModifier(p.SqModifier); err != nil {
		return err
	}
	if err := validateSafetyFactor(p.SafetyFactor); err != nil {
		return err
	}
	if err := validateWhitelistingEnabled(p.WhitelistingEnabled); err != nil {
		return err
	}
	if err := validateInvariantCheckEpoch(p.InvariantCheckEpoch); err != nil {
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
