package types

import "fmt"

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
	if p.EpochIdentifier == "" {
		return fmt.Errorf("epoch_identifier is required")
	}
	return nil
}
