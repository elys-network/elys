package types

import (
	"cosmossdk.io/math"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(poolCreationFee math.Int) Params {
	return Params{
		PoolCreationFee: poolCreationFee,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		math.NewInt(10_000_000), // 10 ELYS
	)
}

// Validate validates the set of params
func (p Params) Validate() error {

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
