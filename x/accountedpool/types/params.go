package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	epochtypes "github.com/elys-network/elys/x/epochs/types"
	"gopkg.in/yaml.v2"
)

var (
	_                      paramtypes.ParamSet = (*Params)(nil)
	KeyInvariantCheckEpoch                     = []byte("InvariantCheckEpoch")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		InvariantCheckEpoch: epochtypes.DayEpochID,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyInvariantCheckEpoch, &p.InvariantCheckEpoch, validateInvariantCheckEpoch),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateInvariantCheckEpoch(p.InvariantCheckEpoch); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateInvariantCheckEpoch(i interface{}) error {
	epoch, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if epoch != epochtypes.DayEpochID && epoch != epochtypes.WeekEpochID && epoch != epochtypes.HourEpochID {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
