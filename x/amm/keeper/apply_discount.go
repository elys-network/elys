package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ApplyDiscount applies discount to swap fee if applicable
func ApplyDiscount(swapFee sdk.Dec, discount sdk.Dec) sdk.Dec {
	// apply discount percentage to swap fee
	swapFee = swapFee.Mul(sdk.OneDec().Sub(discount))
	return swapFee
}
