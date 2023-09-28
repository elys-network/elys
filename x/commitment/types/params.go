package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyVestingInfos = []byte("VestingInfos")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(vestingInfos []*VestingInfo) Params {
	return Params{
		VestingInfos: vestingInfos,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(nil)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVestingInfos, &p.VestingInfos, validateVestingInfos),
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

// Add validators for the new parameters
func validateVestingInfos(i interface{}) error {
	vestingInfos, ok := i.([]*VestingInfo)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, vestingInfo := range vestingInfos {
		if err := validateVestingInfo(vestingInfo); err != nil {
			return err
		}
	}

	return nil
}

func validateVestingInfo(info *VestingInfo) error {
	if info.BaseDenom == "" {
		return fmt.Errorf("base_denom cannot be empty")
	}

	if info.VestingDenom == "" {
		return fmt.Errorf("vesting_denom cannot be empty")
	}

	if info.EpochIdentifier == "" {
		return fmt.Errorf("epoch_identifier cannot be empty")
	}

	if info.NumEpochs <= 0 {
		return fmt.Errorf("num_epochs must be greater than zero")
	}

	if info.NumMaxVestings <= 0 {
		return fmt.Errorf("num_max_vestings must be greater than zero")
	}

	return nil
}
