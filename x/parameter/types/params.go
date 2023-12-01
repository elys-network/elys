package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams(minCommissionRate sdk.Dec, maxVotingPower sdk.Dec, minSelfDelegation sdk.Int, brokerAddress string) Params {
	return Params{
		MinCommissionRate: minCommissionRate,
		MaxVotingPower:    maxVotingPower,
		MinSelfDelegation: minSelfDelegation,
		BrokerAddress:     brokerAddress,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(5, 2),
		sdk.NewDec(100),
		sdk.OneInt(),
		"",
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
