package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcTakeAmount calculates the take amount in the custody asset based on the funding rate
func CalcTakeAmount(custodyAmount sdk.Int, custodyAsset string, fundingRate sdk.Dec) sdk.Int {
	absoluteFundingRate := fundingRate.Abs()

	// Calculate the take amount
	takeAmountValue := sdk.NewDecFromInt(custodyAmount).Mul(absoluteFundingRate).TruncateInt()

	return takeAmountValue
}
