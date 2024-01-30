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
		"elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l",
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMinCommissionRate(p.MinCommissionRate); err != nil {
		return err
	}
	if err := validateMaxVotingPower(p.MaxVotingPower); err != nil {
		return err
	}
	if err := validateMinSelfDelegation(p.MinSelfDelegation); err != nil {
		return err
	}
	if err := validateBrokerAddress(p.BrokerAddress); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMinCommissionRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidMinCommissionRate
	}
	if v.IsNegative() {
		return ErrInvalidMinCommissionRate
	}
	return nil
}

func validateMaxVotingPower(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return ErrInvalidMaxVotingPower
	}
	if v.IsNegative() {
		return ErrInvalidMaxVotingPower
	}
	return nil
}

func validateMinSelfDelegation(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return ErrInvalidMinSelfDelegation
	}
	if v.IsNegative() {
		return ErrInvalidMinSelfDelegation
	}
	return nil
}

func validateBrokerAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return ErrInvalidBrokerAddress
	}
	if v == "" {
		return ErrInvalidBrokerAddress
	}
	return nil
}
