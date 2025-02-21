package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewParams creates a new Params instance
func NewParams(
	minCommissionRate sdkmath.LegacyDec,
	maxVotingPower sdkmath.LegacyDec,
	minSelfDelegation sdkmath.Int,
	totalBlocksPerYear uint64,
	rewardsDataLifeTime uint64,
	takerFees sdkmath.LegacyDec,
	takerFeeAddress string,
) Params {
	return Params{
		MinCommissionRate:         minCommissionRate,
		MaxVotingPower:            maxVotingPower,
		MinSelfDelegation:         minSelfDelegation,
		TotalBlocksPerYear:        totalBlocksPerYear,
		RewardsDataLifetime:       rewardsDataLifeTime,
		TakerFees:                 takerFees,
		TakerFeeCollectionAddress: takerFeeAddress,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		sdkmath.LegacyNewDecWithPrec(5, 2),
		sdkmath.LegacyNewDec(100),
		sdkmath.OneInt(),
		6307200,
		86400, // 1 day
		sdkmath.LegacyZeroDec(),
		authtypes.NewModuleAddress("taker_fee_collection").String(),
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.MinCommissionRate.IsNil() {
		return fmt.Errorf("minimum commission rate cannot be nil")
	}
	if p.MinCommissionRate.IsNegative() {
		return ErrInvalidMinCommissionRate
	}

	if p.MaxVotingPower.IsNil() {
		return fmt.Errorf("maximum voting power cannot be nil")
	}
	if p.MaxVotingPower.IsNegative() {
		return ErrInvalidMaxVotingPower
	}

	if p.MinCommissionRate.IsNil() {
		return fmt.Errorf("minimum commission rate cannot be nil")
	}
	if p.MinCommissionRate.IsNegative() {
		return ErrInvalidMinSelfDelegation
	}

	if p.TotalBlocksPerYear <= 0 {
		return fmt.Errorf("total blocks per year cannot be negative or zero")
	}

	if p.RewardsDataLifetime <= 0 {
		return fmt.Errorf("rewards data lifetime cannot be negative or zero")
	}
	return nil
}
