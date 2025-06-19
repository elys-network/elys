package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		VestingInfos: []VestingInfo{
			{
				BaseDenom:      ptypes.Eden,
				VestingDenom:   ptypes.Elys,
				NumBlocks:      1576800,
				VestNowFactor:  math.NewInt(90),
				NumMaxVestings: 10000,
			},
		},
		TotalCommitted: sdk.Coins{},
		EnableVestNow:  false,
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, vestingInfo := range p.VestingInfos {
		if err := vestingInfo.Validate(); err != nil {
			return err
		}
	}
	if err := p.TotalCommitted.Validate(); err != nil {
		return err
	}
	return nil
}

func (vestingInfo VestingInfo) Validate() error {
	if err := sdk.ValidateDenom(vestingInfo.BaseDenom); err != nil {
		return errorsmod.Wrapf(ErrInvalidDenom, vestingInfo.BaseDenom)
	}
	if err := sdk.ValidateDenom(vestingInfo.VestingDenom); err != nil {
		return errorsmod.Wrapf(ErrInvalidDenom, vestingInfo.VestingDenom)
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
	return nil
}
