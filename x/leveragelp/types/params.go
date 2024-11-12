package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"

	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		LeverageMax:         sdkmath.LegacyNewDec(10),
		EpochLength:         (int64)(1),
		MaxOpenPositions:    (int64)(9999),
		PoolOpenThreshold:   sdkmath.LegacyNewDecWithPrec(2, 1),  // 0.2
		SafetyFactor:        sdkmath.LegacyNewDecWithPrec(11, 1), // 1.1
		WhitelistingEnabled: false,
		FallbackEnabled:     true,
		NumberPerBlock:      (int64)(1000),
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
	if err := validateEpochLength(p.EpochLength); err != nil {
		return err
	}
	if err := validateMaxOpenPositions(p.MaxOpenPositions); err != nil {
		return err
	}
	if err := validatePoolOpenThreshold(p.PoolOpenThreshold); err != nil {
		return err
	}
	if err := validateSafetyFactor(p.SafetyFactor); err != nil {
		return err
	}
	if err := validateWhitelistingEnabled(p.WhitelistingEnabled); err != nil {
		return err
	}
	if err := validateNumberOfBlocks(p.NumberPerBlock); err != nil {
		return err
	}
	if err := validateFallbackEnabled(p.FallbackEnabled); err != nil {
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
	if !v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("leverage max must be greater than 1: %s", v)
	}
	if v.GT(sdkmath.LegacyNewDec(10)) {
		return fmt.Errorf("leverage max too large: %s", v)
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

func validateSafetyFactor(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("safety factor must be not nil")
	}
	if !v.IsPositive() {
		return fmt.Errorf("safety factor must be positive: %s", v)
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
	if !v.IsPositive() {
		return fmt.Errorf("pool open threshold must be positive: %s", v)
	}

	return nil
}

func validateNumberOfBlocks(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("number of positions per block must be positive: %d", v)
	}

	if v > MaxPageLimit {
		return fmt.Errorf("number of positions per block should not exceed page limit: %d, number of positions: %d", MaxPageLimit, v)
	}

	return nil
}

func validateFallbackEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
