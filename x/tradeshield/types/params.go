package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

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
		RewardPercentage:     sdk.ZeroDec(),
		MarginError:          sdk.ZeroDec(),
		MinimumDeposit:       sdk.ZeroInt(),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
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
