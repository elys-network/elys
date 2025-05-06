package types

import (
	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcTakeAmount calculates the take amount in the custody asset based on the funding rate
func CalcTakeAmount(custodyAmount sdkmath.Int, fundingRate osmomath.BigDec) sdkmath.Int {
	absoluteFundingRate := fundingRate.Abs()

	// Calculate the take amount
	takeAmountValue := osmomath.BigDecFromSDKInt(custodyAmount).Mul(absoluteFundingRate).Dec().TruncateInt()

	return takeAmountValue
}
