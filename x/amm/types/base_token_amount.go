package types

import (
	sdkmath "cosmossdk.io/math"
)

func BaseTokenAmount(decimal uint64) sdkmath.Int {
	return sdkmath.NewIntWithDecimal(1, int(decimal))
}
