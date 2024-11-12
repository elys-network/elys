package types

import (
	sdkmath "cosmossdk.io/math"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MarketOrderEnabled:   true,
		StakeEnabled:         true,
		ProcessOrdersEnabled: true,
		SwapEnabled:          true,
		PerpetualEnabled:     true,
		RewardEnabled:        true,
		LeverageEnabled:      true,
		LimitProcessOrder:    1000000,
		RewardPercentage:     sdkmath.LegacyZeroDec(),
		MarginError:          sdkmath.LegacyZeroDec(),
		MinimumDeposit:       sdkmath.ZeroInt(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
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
