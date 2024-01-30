package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalculateFundingRate calculates and returns the funding rate based on long and short amounts
func CalcFundingRate(longAmount, shortAmount sdk.Int, baseRate, maxRate, minRate sdk.Dec) sdk.Dec {
	var ratio sdk.Dec

	// Check for division by zero when longAmount > shortAmount
	if longAmount.GT(shortAmount) {
		if shortAmount.IsZero() {
			// Handle the case where shortAmount is zero
			return maxRate
		}
		ratio = sdk.NewDecFromInt(longAmount).Quo(sdk.NewDecFromInt(shortAmount))
		return sdk.MinDec(sdk.MaxDec(baseRate.Mul(ratio), minRate), maxRate)
	} else if shortAmount.GT(longAmount) {
		if longAmount.IsZero() {
			// Handle the case where longAmount is zero
			return maxRate
		}
		ratio = sdk.NewDecFromInt(shortAmount).Quo(sdk.NewDecFromInt(longAmount))
		return sdk.MinDec(sdk.MaxDec(baseRate.Mul(ratio).Neg(), minRate), maxRate)
	} else {
		// In case of exact equality, return the base rate
		return baseRate
	}
}
