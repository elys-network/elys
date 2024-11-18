package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		BorrowInterestRateDecrease:                     math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateIncrease:                     math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateMax:                          math.LegacyMustNewDecFromStr("0.3"),
		BorrowInterestRateMin:                          math.LegacyMustNewDecFromStr("0.1"),
		EnableTakeProfitCustodyLiabilities:             false,
		FixedFundingRate:                               math.LegacyMustNewDecFromStr("0.5"), // 50%
		ForceCloseFundAddress:                          authtypes.NewModuleAddress("zero").String(),
		ForceCloseFundPercentage:                       math.LegacyOneDec(),
		HealthGainFactor:                               math.LegacyMustNewDecFromStr("0.000000220000000000"),
		IncrementalBorrowInterestPaymentEnabled:        true,
		IncrementalBorrowInterestPaymentFundAddress:    authtypes.NewModuleAddress("zero").String(),
		IncrementalBorrowInterestPaymentFundPercentage: math.LegacyMustNewDecFromStr("0.1"),
		LeverageMax:                                    math.LegacyNewDec(25),
		MaxLimitOrder:                                  (int64)(100000),
		MaxOpenPositions:                               (int64)(100000),
		MaximumLongTakeProfitPriceRatio:                math.LegacyMustNewDecFromStr("11"),
		MaximumShortTakeProfitPriceRatio:               math.LegacyMustNewDecFromStr("0.98"),
		MinimumLongTakeProfitPriceRatio:                math.LegacyMustNewDecFromStr("1.02"),
		PerpetualSwapFee:                               math.LegacyMustNewDecFromStr("0.001"), // 0.1%
		PoolOpenThreshold:                              math.LegacyMustNewDecFromStr("0.65"),
		SafetyFactor:                                   math.LegacyMustNewDecFromStr("1.025000000000000000"),
		WeightBreakingFeeFactor:                        math.LegacyMustNewDecFromStr("0.5"),
		WhitelistingEnabled:                            false,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

func CheckLegacyDecNilAndNegative(value math.LegacyDec, name string) error {
	if value.IsNil() {
		return fmt.Errorf("%s is nil", name)
	}
	if value.IsNegative() {
		return fmt.Errorf("%s is negative", name)
	}
	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := CheckLegacyDecNilAndNegative(p.BorrowInterestRateMax, "BorrowInterestRateMax"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.BorrowInterestRateMin, "BorrowInterestRateMin"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.BorrowInterestRateIncrease, "BorrowInterestRateIncrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.BorrowInterestRateDecrease, "BorrowInterestRateDecrease"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.FixedFundingRate, "FixedFundingRate"); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.ForceCloseFundAddress); err != nil {
		return fmt.Errorf("invalid address: %s", err.Error())
	}
	if err := CheckLegacyDecNilAndNegative(p.ForceCloseFundPercentage, "ForceCloseFundPercentage"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.HealthGainFactor, "HealthGainFactor"); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.IncrementalBorrowInterestPaymentFundAddress); err != nil {
		return fmt.Errorf("invalid address: %s", err.Error())
	}
	if err := CheckLegacyDecNilAndNegative(p.IncrementalBorrowInterestPaymentFundPercentage, "IncrementalBorrowInterestPaymentFundPercentage"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LeverageMax, "LeverageMax"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.MaximumLongTakeProfitPriceRatio, "MaximumLongTakeProfitPriceRatio"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.MaximumShortTakeProfitPriceRatio, "MaximumShortTakeProfitPriceRatio"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.MinimumLongTakeProfitPriceRatio, "MinimumLongTakeProfitPriceRatio"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.PerpetualSwapFee, "PerpetualSwapFee"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.PoolOpenThreshold, "PoolOpenThreshold"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.SafetyFactor, "SafetyFactor"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.WeightBreakingFeeFactor, "WeightBreakingFeeFactor"); err != nil {
		return err
	}

	return nil
}
