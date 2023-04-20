package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyEpochIdentifier = []byte("EpochIdentifier")
	// TODO: Determine the default value
	DefaultEpochIdentifier string = "epoch_identifier"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

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

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, validateEpochIdentifier),
	}
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
