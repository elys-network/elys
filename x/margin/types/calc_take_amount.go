package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcTakeAmount calculates the take amount in the custody asset based on the funding rate
func CalcTakeAmount(custodyAmount sdk.Coin, custodyAsset string, fundingRate sdk.Dec) sdk.Coin {
	absoluteFundingRate := fundingRate.Abs()

	// Calculate the take amount
	takeAmountValue := sdk.NewDecFromInt(custodyAmount.Amount).Mul(absoluteFundingRate).TruncateInt()

	// Create the take amount coin
	takeAmount := sdk.NewCoin(custodyAsset, takeAmountValue)

	return takeAmount
}
