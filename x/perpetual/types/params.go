package types

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		BorrowInterestRateDecrease:          math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateIncrease:          math.LegacyMustNewDecFromStr("0.0003"),
		BorrowInterestRateMax:               math.LegacyMustNewDecFromStr("0.3"),
		BorrowInterestRateMin:               math.LegacyMustNewDecFromStr("0.1"),
		FixedFundingRate:                    math.LegacyMustNewDecFromStr("0.5"), // 50%
		HealthGainFactor:                    math.LegacyMustNewDecFromStr("0.000000220000000000"),
		BorrowInterestPaymentEnabled:        true,
		BorrowInterestPaymentFundPercentage: math.LegacyMustNewDecFromStr("0.1"),
		LeverageMax:                         math.LegacyNewDec(25),
		MaxLimitOrder:                       (int64)(100000),
		MaxOpenPositions:                    (int64)(100000),
		MaximumLongTakeProfitPriceRatio:     math.LegacyMustNewDecFromStr("11"),
		MaximumShortTakeProfitPriceRatio:    math.LegacyMustNewDecFromStr("0.98"),
		MinimumLongTakeProfitPriceRatio:     math.LegacyMustNewDecFromStr("1.02"),
		PerpetualSwapFee:                    math.LegacyMustNewDecFromStr("0.001"), // 0.1%
		PoolMaxLiabilitiesThreshold:         math.LegacyMustNewDecFromStr("0.65"),
		SafetyFactor:                        math.LegacyMustNewDecFromStr("1.025000000000000000"),
		WeightBreakingFeeFactor:             math.LegacyMustNewDecFromStr("0.5"),
		WhitelistingEnabled:                 false,
		EnabledPools:                        []uint64(nil),
		MinimumNotionalValue:                math.LegacyNewDec(10), // Minimum position notional value of 10 units
		LongMinimumLiabilityAmount:          math.NewInt(1),
		ExitBuffer:                          math.LegacyMustNewDecFromStr("0.15"),
		TakerFee:                            math.LegacyMustNewDecFromStr("0.00075"),
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
		return errors.New("BorrowInterestRateMin must be less than BorrowInterestRateMax")
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
		return errors.New("MaximumLongTakeProfitPriceRatio must be greater than 1")
	}
	if err := CheckLegacyDecNilAndNegative(p.MaximumShortTakeProfitPriceRatio, "MaximumShortTakeProfitPriceRatio"); err != nil {
		return err
	}
	if p.MaximumShortTakeProfitPriceRatio.GTE(math.LegacyOneDec()) {
		return errors.New("MaximumShortTakeProfitPriceRatio must be less than 1")
	}
	if err := CheckLegacyDecNilAndNegative(p.MinimumLongTakeProfitPriceRatio, "MinimumLongTakeProfitPriceRatio"); err != nil {
		return err
	}
	if p.MinimumLongTakeProfitPriceRatio.LTE(math.LegacyOneDec()) {
		return errors.New("MinimumLongTakeProfitPriceRatio must be greater than 1")
	}
	if p.MaximumLongTakeProfitPriceRatio.LTE(p.MinimumLongTakeProfitPriceRatio) {
		return errors.New("MaximumLongTakeProfitPriceRatio must be greater than MinimumLongTakeProfitPriceRatio")
	}
	if err := CheckLegacyDecNilAndNegative(p.PerpetualSwapFee, "PerpetualSwapFee"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.PoolMaxLiabilitiesThreshold, "PoolMaxLiabilitiesThreshold"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.SafetyFactor, "SafetyFactor"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.WeightBreakingFeeFactor, "WeightBreakingFeeFactor"); err != nil {
		return err
	}

	if containsDuplicates(p.EnabledPools) {
		return errors.New("array must not contain duplicate values")
	}

	if err := CheckLegacyDecNilAndNegative(p.ExitBuffer, "ExitBuffer"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.TakerFee, "TakerFee"); err != nil {
		return err
	}

	if p.TakerFee.GTE(math.LegacyOneDec()) {
		return errors.New("TakerFee must be less than 1")
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

func (p Params) GetBigDecFixedFundingRate() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.FixedFundingRate)
}

func (p Params) GetBigDecPerpetualSwapFee() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.PerpetualSwapFee)
}

func (p Params) GetBigDecWeightBreakingFeeFactor() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.WeightBreakingFeeFactor)
}

func (p Params) GetBigDecMinimumLongTakeProfitPriceRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MinimumLongTakeProfitPriceRatio)
}

func (p Params) GetBigDecMaximumLongTakeProfitPriceRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaximumLongTakeProfitPriceRatio)
}

func (p Params) GetBigDecMaximumShortTakeProfitPriceRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaximumShortTakeProfitPriceRatio)
}

func (p Params) GetBigDecSafetyFactor() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.SafetyFactor)
}

func (p Params) GetBigDecBorrowInterestRateMin() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.BorrowInterestRateMin)
}

func (p Params) GetBigDecBorrowInterestPaymentFundPercentage() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.BorrowInterestPaymentFundPercentage)
}

func (p Params) GetBigDecTakerFees() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.TakerFee)
}
