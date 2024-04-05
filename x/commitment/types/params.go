package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*LegacyParams)(nil)

var KeyVestingInfos = []byte("VestingInfos")

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&LegacyParams{})
}

// NewLegacyParams creates a new LegacyParams instance
func NewLegacyParams(vestingInfos []*LegacyVestingInfo) LegacyParams {
	return LegacyParams{
		VestingInfos: vestingInfos,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		VestingInfos:   nil,
		TotalCommitted: sdk.Coins(nil),
	}
}

// ParamSetPairs get the params.ParamSet
func (p *LegacyParams) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVestingInfos, &p.VestingInfos, validateVestingInfos),
	}
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

	if info.NumBlocks <= 0 {
		return fmt.Errorf("num_blocks must be greater than zero")
	}

	if info.NumMaxVestings < 0 {
		return fmt.Errorf("num_max_vestings cannot be negative")
	}

	return nil
}
