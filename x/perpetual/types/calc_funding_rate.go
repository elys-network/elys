package types

import (
	sdkmath "cosmossdk.io/math"
)

// CalculateFundingRate calculates and returns the funding rate based on long and short amounts
func CalcFundingRate(longAmount, shortAmount sdkmath.Int, baseRate, maxRate, minRate sdkmath.LegacyDec) sdkmath.LegacyDec {
	var ratio sdkmath.LegacyDec

	// Check for division by zero when longAmount > shortAmount
	if longAmount.GT(shortAmount) {
		if shortAmount.IsZero() {
			// Handle the case where shortAmount is zero
			return maxRate
		}
		ratio = longAmount.ToLegacyDec().Quo(shortAmount.ToLegacyDec())
		return sdkmath.LegacyMinDec(sdkmath.LegacyMaxDec(baseRate.Mul(ratio), minRate), maxRate)
	} else if shortAmount.GT(longAmount) {
		if longAmount.IsZero() {
			// Handle the case where longAmount is zero
			return maxRate
		}
		ratio = shortAmount.ToLegacyDec().Quo(longAmount.ToLegacyDec())
		return sdkmath.LegacyMinDec(sdkmath.LegacyMaxDec(baseRate.Mul(ratio).Neg(), minRate), maxRate)
	} else {
		// In case of exact equality, return the base rate
		return baseRate
	}
}
