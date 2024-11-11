package types

import sdkmath "cosmossdk.io/math"

// CalcTakeAmount calculates the take amount in the custody asset based on the funding rate
func CalcTakeAmount(custodyAmount sdkmath.Int, fundingRate sdkmath.LegacyDec) sdkmath.Int {
	absoluteFundingRate := fundingRate.Abs()

	// Calculate the take amount
	takeAmountValue := custodyAmount.ToLegacyDec().Mul(absoluteFundingRate).TruncateInt()

	return takeAmountValue
}
