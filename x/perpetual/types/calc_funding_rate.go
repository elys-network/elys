package types

import (
	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcFundingRate calculates and returns the funding rate based on long and short amounts
func CalcFundingRate(longAmount, shortAmount sdkmath.Int, baseRate, maxRate, minRate osmomath.BigDec) osmomath.BigDec {
	var ratio osmomath.BigDec

	// Check for division by zero when longAmount > shortAmount
	if longAmount.GT(shortAmount) {
		if shortAmount.IsZero() {
			// Handle the case where shortAmount is zero
			return maxRate
		}
		ratio = osmomath.BigDecFromSDKInt(longAmount).Quo(osmomath.BigDecFromSDKInt(shortAmount))
		return osmomath.MinBigDec(osmomath.MaxBigDec(baseRate.Mul(ratio), minRate), maxRate)
	} else if shortAmount.GT(longAmount) {
		if longAmount.IsZero() {
			// Handle the case where longAmount is zero
			return maxRate
		}
		ratio = osmomath.BigDecFromSDKInt(shortAmount).Quo(osmomath.BigDecFromSDKInt(longAmount))
		return osmomath.MinBigDec(osmomath.MaxBigDec(baseRate.Mul(ratio).Neg(), minRate), maxRate)
	} else {
		// In case of exact equality, return the base rate
		return baseRate
	}
}
