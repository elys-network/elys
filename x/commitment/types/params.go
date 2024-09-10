package types

import (
	"fmt"
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
	for _, vestingInfo := range p.VestingInfos {
		if vestingInfo == nil {
			return fmt.Errorf("vesting info cannot be nil")
		}
		if err := sdk.ValidateDenom(vestingInfo.BaseDenom); err != nil {
			return err
		}
		if err := sdk.ValidateDenom(vestingInfo.VestingDenom); err != nil {
			return err
		}
		if vestingInfo.NumMaxVestings < 0 {
			return fmt.Errorf("num_max_vestings cannot be negative")
		}
		if vestingInfo.NumBlocks < 0 {
			return fmt.Errorf("num_blocks cannot be negative")
		}
		if vestingInfo.VestNowFactor.IsNil() {
			return fmt.Errorf("vesting now factor cannot be nil")
		}
		if !vestingInfo.VestNowFactor.IsPositive() {
			return fmt.Errorf("vesting now factor must be positive")
		}
	}
	if err := p.TotalCommitted.Validate(); err != nil {
		return err
	}
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
