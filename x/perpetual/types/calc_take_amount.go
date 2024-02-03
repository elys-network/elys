package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcTakeAmount calculates the take amount in the custody asset based on the funding rate
func CalcTakeAmount(custodyAmount math.Int, custodyAsset string, fundingRate sdk.Dec) math.Int {
	absoluteFundingRate := fundingRate.Abs()

	// Calculate the take amount
	takeAmountValue := sdk.NewDecFromInt(custodyAmount).Mul(absoluteFundingRate).TruncateInt()

	return takeAmountValue
}
