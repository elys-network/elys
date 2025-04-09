package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MarketOrderEnabled:   true,
		StakeEnabled:         true,
		ProcessOrdersEnabled: true,
		SwapEnabled:          true,
		PerpetualEnabled:     true,
		RewardEnabled:        true,
		LeverageEnabled:      true,
		LimitProcessOrder:    1000000,
		RewardPercentage:     sdkmath.LegacyZeroDec(),
		MarginError:          sdkmath.LegacyZeroDec(),
		MinimumDeposit:       sdkmath.ZeroInt(),
		Tolerance:            sdkmath.LegacyMustNewDecFromStr("0.05"),
	}
}

func CheckLegacyDecNilAndNegative(value sdkmath.LegacyDec, name string) error {
	if value.IsNil() {
		return fmt.Errorf("%s is nil", name)
	}
	if value.IsNegative() {
		return fmt.Errorf("%s is negative", name)
	}
	return nil
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := CheckLegacyDecNilAndNegative(p.RewardPercentage, "RewardPercentage"); err != nil {
		return err
	}
	if err := CheckLegacyDecNilAndNegative(p.MarginError, "MarginError"); err != nil {
		return err
	}
	if p.MinimumDeposit.IsNil() {
		return fmt.Errorf("MinimumDeposit is required")
	}
	if p.MinimumDeposit.IsNegative() {
		return fmt.Errorf("MinimumDeposit is negative")
	}
	return nil
}

func (p Params) GetBigDecTolerance() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.Tolerance)
}
