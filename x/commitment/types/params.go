package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		VestingInfos:   nil,
		TotalCommitted: sdk.Coins(nil),
	}
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

// Validate validates the set of params
func (p LegacyParams) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p LegacyParams) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
