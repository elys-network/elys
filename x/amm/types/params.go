package types

import (
	"cosmossdk.io/math"
	"fmt"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(poolCreationFee math.Int, slippageTrackDuration uint64, enable bool) Params {
	return Params{
		PoolCreationFee:       poolCreationFee,
		SlippageTrackDuration: slippageTrackDuration,
		EnableUsdcPairedPoolOnly: enable,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		math.NewInt(10_000_000), // 10 ELYS
		86400*7,
		false,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.PoolCreationFee.IsNil() {
		return fmt.Errorf("pool creation fee must not be empty")
	}
	if p.PoolCreationFee.IsNegative() {
		return fmt.Errorf("pool creation fee must be positive")
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
