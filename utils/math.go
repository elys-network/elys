package utils

import (
	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func Pow10(decimal uint64) osmomath.BigDec {
	if decimal <= 18 {
		result := int64(1)
		for i := int64(0); i < int64(decimal); i++ {
			result = result * 10
		}
		return osmomath.NewBigDec(result)
	}

	// This case less likely to happen
	value := osmomath.NewBigDec(1)
	for i := 0; i < int(decimal); i++ {
		value = value.MulInt64(10)
	}
	return value
}

// AbsDifferenceWithSign returns | a - b |, (a - b).sign()
// a is mutated and returned.
func AbsDifferenceWithSign(a, b sdkmath.LegacyDec) (sdkmath.LegacyDec, bool) {
	if a.GTE(b) {
		return a.SubMut(b), false
	} else {
		return a.NegMut().AddMut(b), true
	}
}
