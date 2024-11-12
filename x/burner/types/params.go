package types

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

const (
	// TODO: Determine the default value
	DefaultEpochIdentifier string = "epoch_identifier"
)

// NewParams creates a new Params instance
func NewParams(
	epochIdentifier string,
) Params {
	return Params{
		EpochIdentifier: epochIdentifier,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultEpochIdentifier,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateEpochIdentifier(p.EpochIdentifier); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateEpochIdentifier validates the EpochIdentifier param
func validateEpochIdentifier(v interface{}) error {
	epochIdentifier, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = epochIdentifier

	return nil
}
