package types

import (
	sdkmath "cosmossdk.io/math"
)

func OneTokenUnit(decimal uint64) sdkmath.Int {
	return sdkmath.NewIntWithDecimal(1, int(decimal))
}
