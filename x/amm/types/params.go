package types

import (
	"errors"
	"slices"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// NewParams creates a new Params instance
func NewParams(poolCreationFee math.Int, slippageTrackDuration uint64, baseAssets []string) Params {
	return Params{
		PoolCreationFee:                  poolCreationFee,
		SlippageTrackDuration:            slippageTrackDuration,
		BaseAssets:                       baseAssets,
		WeightBreakingFeeExponent:        math.LegacyMustNewDecFromStr("2.5"),
		WeightBreakingFeeMultiplier:      math.LegacyMustNewDecFromStr("0.0005"),
		WeightBreakingFeePortion:         math.LegacyMustNewDecFromStr("0.5"),
		WeightRecoveryFeePortion:         math.LegacyMustNewDecFromStr("0.1"),
		ThresholdWeightDifference:        math.LegacyMustNewDecFromStr("0.04"),
		AllowedPoolCreators:              []string{authtypes.NewModuleAddress(govtypes.ModuleName).String()},
		ThresholdWeightDifferenceSwapFee: math.LegacyMustNewDecFromStr("0.15"),
		LpLockupDuration:                 60,
		MinSlippage:                      math.LegacyMustNewDecFromStr("0.0005"),
		AllowedUpfrontSwapMakers:         []string{},
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
		return errors.New("pool creation fee must be positive")
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
	if p.WeightBreakingFeeExponent.GT(math.LegacyMustNewDecFromStr("3")) {
		return errors.New("weightBreakingFeeExponent must be less than 3")
	}

	if p.WeightBreakingFeeMultiplier.IsNil() {
		return errors.New("weightBreakingFeeMultiplier must not be empty")
	}
	if p.WeightBreakingFeeMultiplier.IsNegative() {
		return errors.New("weightBreakingFeeMultiplier must be positive")
	}
	if p.WeightBreakingFeeMultiplier.GT(math.LegacyMustNewDecFromStr("0.001")) {
		return errors.New("weightBreakingFeeMultiplier must be less than 0.01%")
	}

	if p.WeightBreakingFeePortion.IsNil() {
		return errors.New("weightBreakingFeePortion must not be empty")
	}
	if p.WeightBreakingFeePortion.IsNegative() {
		return errors.New("weightBreakingFeePortion must be positive")
	}
	if p.WeightBreakingFeePortion.GT(math.LegacyMustNewDecFromStr("1")) {
		return errors.New("weightBreakingFeePortion must be less than 1")
	}

	if p.WeightRecoveryFeePortion.IsNil() {
		return errors.New("weightRecoveryFeePortion must not be empty")
	}
	if p.WeightRecoveryFeePortion.IsNegative() {
		return errors.New("weightRecoveryFeePortion must be positive")
	}
	if p.WeightRecoveryFeePortion.GT(math.LegacyMustNewDecFromStr("1")) {
		return errors.New("weightRecoveryFeePortion must be less than 1")
	}

	if p.ThresholdWeightDifference.IsNil() {
		return errors.New("thresholdWeightDifference must not be empty")
	}
	if p.ThresholdWeightDifference.IsNegative() {
		return errors.New("thresholdWeightDifference must be positive")
	}
	if p.ThresholdWeightDifference.GT(math.LegacyMustNewDecFromStr("0.1")) {
		return errors.New("thresholdWeightDifference must be less than 0.1%")
	}

	if p.MinSlippage.IsNil() {
		return errors.New("MinSlippage must not be empty")
	}
	if p.MinSlippage.IsNegative() {
		return errors.New("MinSlippage must be positive")
	}
	if p.MinSlippage.GT(math.LegacyMustNewDecFromStr("0.01")) {
		return errors.New("MinSlippage must be less than 1%")
	}

	if p.ThresholdWeightDifferenceSwapFee.IsNil() {
		return errors.New("thresholdWeightDifferenceSwapFee must not be empty")
	}
	if p.ThresholdWeightDifferenceSwapFee.IsNegative() {
		return errors.New("thresholdWeightDifferenceSwapFee must be positive")
	}
	if p.ThresholdWeightDifferenceSwapFee.GT(math.LegacyMustNewDecFromStr("0.15")) {
		return errors.New("thresholdWeightDifferenceSwapFee must be less than 0.15%")
	}

	for _, swapMaker := range p.AllowedUpfrontSwapMakers {
		if _, err := sdk.AccAddressFromBech32(swapMaker); err != nil {
			return errors.New("invalid upfront swap maker address")
		}
	}

	return nil
}

func (p Params) IsCreatorAllowed(creator string) bool {
	return slices.Contains(p.AllowedPoolCreators, creator)
}

func (p Params) GetBigDecMinSlippage() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MinSlippage)
}

func (p Params) GetBigDecWeightBreakingFeeMultiplier() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.WeightBreakingFeeMultiplier)
}

func (p Params) GetBigDecWeightBreakingFeeExponent() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.WeightBreakingFeeExponent)
}

func (p Params) GetBigDecThresholdWeightDifference() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.ThresholdWeightDifference)
}

func (p Params) GetBigDecThresholdWeightDifferenceSwapFee() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.ThresholdWeightDifferenceSwapFee)
}

func (p Params) GetBigDecWeightBreakingFeePortion() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.WeightBreakingFeePortion)
}

func (p Params) GetBigDecWeightRecoveryFeePortion() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.WeightRecoveryFeePortion)
}
