package types

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewParams creates a new Params instance
func NewParams(poolCreationFee math.Int, slippageTrackDuration uint64, baseAssets []string) Params {
	return Params{
		PoolCreationFee:             poolCreationFee,
		SlippageTrackDuration:       slippageTrackDuration,
		BaseAssets:                  baseAssets,
		WeightBreakingFeeExponent:   math.LegacyMustNewDecFromStr("2.5"),
		WeightBreakingFeeMultiplier: math.LegacyMustNewDecFromStr("0.0005"),
		WeightBreakingFeePortion:    math.LegacyMustNewDecFromStr("0.5"),
		WeightRecoveryFeePortion:    math.LegacyMustNewDecFromStr("0.1"),
		ThresholdWeightDifference:   math.LegacyMustNewDecFromStr("0.3"),
		AllowedPoolCreators:         []string{authtypes.NewModuleAddress(govtypes.ModuleName).String()},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		math.NewInt(10_000_000), // 10 ELYS
		86400*7,
		[]string{},
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.PoolCreationFee.IsNil() {
		return errors.New("pool creation fee must not be empty")
	}
	if p.PoolCreationFee.IsNegative() {
		return fmt.Errorf("pool creation fee must be positive")
	}

	for _, asset := range p.BaseAssets {
		if err := sdk.ValidateDenom(asset); err != nil {
			return err
		}
	}

	if p.WeightBreakingFeeExponent.IsNil() {
		return errors.New("weightBreakingFeeExponent must not be empty")
	}
	if p.WeightBreakingFeeExponent.IsNegative() {
		return errors.New("weightBreakingFeeExponent must be positive")
	}

	if p.WeightBreakingFeeMultiplier.IsNil() {
		return errors.New("weightBreakingFeeMultiplier must not be empty")
	}
	if p.WeightBreakingFeeMultiplier.IsNegative() {
		return errors.New("weightBreakingFeeMultiplier must be positive")
	}

	if p.WeightBreakingFeePortion.IsNil() {
		return errors.New("weightBreakingFeePortion must not be empty")
	}
	if p.WeightBreakingFeePortion.IsNegative() {
		return errors.New("weightBreakingFeePortion must be positive")
	}

	if p.WeightRecoveryFeePortion.IsNil() {
		return errors.New("weightRecoveryFeePortion must not be empty")
	}
	if p.WeightRecoveryFeePortion.IsNegative() {
		return errors.New("weightRecoveryFeePortion must be positive")
	}

	if p.ThresholdWeightDifference.IsNil() {
		return errors.New("thresholdWeightDifference must not be empty")
	}
	if p.ThresholdWeightDifference.IsNegative() {
		return errors.New("thresholdWeightDifference must be positive")
	}
	return nil
}

func (p Params) IsCreatorAllowed(creator string) bool {
	for _, allowedCreator := range p.AllowedPoolCreators {
		if allowedCreator == creator {
			return true
		}
	}
	return false
}
