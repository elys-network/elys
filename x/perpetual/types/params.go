package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		BorrowInterestRateDecrease:          math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateIncrease:          math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateMax:               math.LegacyMustNewDecFromStr("0.3"),
		BorrowInterestRateMin:               math.LegacyMustNewDecFromStr("0.1"),
		EnableTakeProfitCustodyLiabilities:  false,
		FixedFundingRate:                    math.LegacyMustNewDecFromStr("0.5"), // 50%
		HealthGainFactor:                    math.LegacyMustNewDecFromStr("0.000000220000000000"),
		BorrowInterestPaymentEnabled:        true,
		BorrowInterestPaymentFundAddress:    authtypes.NewModuleAddress("borrow-interest-payment-fund").String(), // IMP: Shouldn't be same as perpetual module address
		BorrowInterestPaymentFundPercentage: math.LegacyMustNewDecFromStr("0.1"),
		LeverageMax:                         math.LegacyNewDec(25),
		MaxLimitOrder:                       (int64)(100000),
		MaxOpenPositions:                    (int64)(100000),
		MaximumLongTakeProfitPriceRatio:     math.LegacyMustNewDecFromStr("11"),
		MaximumShortTakeProfitPriceRatio:    math.LegacyMustNewDecFromStr("0.98"),
		MinimumLongTakeProfitPriceRatio:     math.LegacyMustNewDecFromStr("1.02"),
		PerpetualSwapFee:                    math.LegacyMustNewDecFromStr("0.001"), // 0.1%
		PoolOpenThreshold:                   math.LegacyMustNewDecFromStr("0.65"),
		SafetyFactor:                        math.LegacyMustNewDecFromStr("1.025000000000000000"),
		WeightBreakingFeeFactor:             math.LegacyMustNewDecFromStr("0.5"),
		WhitelistingEnabled:                 false,
		EnabledPools:                        []uint64{},
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
	if p.BorrowInterestRateMin.GT(p.BorrowInterestRateMax) {
		return fmt.Errorf("BorrowInterestRateMin must be less than BorrowInterestRateMax")
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
	if err := CheckLegacyDecNilAndNegative(p.HealthGainFactor, "HealthGainFactor"); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.BorrowInterestPaymentFundAddress); err != nil {
		return fmt.Errorf("invalid address: %s", err.Error())
	}
	if err := CheckLegacyDecNilAndNegative(p.BorrowInterestPaymentFundPercentage, "IncrementalBorrowInterestPaymentFundPercentage"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.LeverageMax, "LeverageMax"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.MaximumLongTakeProfitPriceRatio, "MaximumLongTakeProfitPriceRatio"); err != nil {
		return err
	}
	if p.MaximumLongTakeProfitPriceRatio.LTE(math.LegacyOneDec()) {
		return fmt.Errorf("MaximumLongTakeProfitPriceRatio must be greater than 1")
	}
	if err := CheckLegacyDecNilAndNegative(p.MaximumShortTakeProfitPriceRatio, "MaximumShortTakeProfitPriceRatio"); err != nil {
		return err
	}
	if p.MaximumShortTakeProfitPriceRatio.GTE(math.LegacyOneDec()) {
		return fmt.Errorf("MaximumShortTakeProfitPriceRatio must be less than 1")
	}
	if err := CheckLegacyDecNilAndNegative(p.MinimumLongTakeProfitPriceRatio, "MinimumLongTakeProfitPriceRatio"); err != nil {
		return err
	}
	if p.MinimumLongTakeProfitPriceRatio.LTE(math.LegacyOneDec()) {
		return fmt.Errorf("MinimumLongTakeProfitPriceRatio must be greater than 1")
	}
	if p.MaximumLongTakeProfitPriceRatio.LTE(p.MinimumLongTakeProfitPriceRatio) {
		return fmt.Errorf("MaximumLongTakeProfitPriceRatio must be greater than MinimumLongTakeProfitPriceRatio")
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

	if containsDuplicates(p.EnabledPools) {
		return fmt.Errorf("array must not contain duplicate values")
	}
	return nil
}

func containsDuplicates(arr []uint64) bool {
	valueMap := make(map[uint64]struct{})
	for _, num := range arr {
		if _, exists := valueMap[num]; exists {
			return true
		}
		valueMap[num] = struct{}{}
	}
	return false
}
