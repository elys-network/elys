package utils

import (
	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

var pow10Int64Cache = [...]int64{
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
}

func Pow10(decimal uint64) osmomath.BigDec {
	if decimal <= 18 {
		return osmomath.NewBigDec(pow10Int64Cache[decimal])
	}

	// This case less likely to happen
	value := osmomath.NewBigDec(1)
	for i := 0; i < int(decimal); i++ {
		value = value.MulInt64(10)
	}
	return value
}

func Pow10Int64(decimal uint64) int64 {
	if decimal <= 18 {
		return pow10Int64Cache[decimal]
	} else {
		panic("cannot do more than 10^18 for int64")
	}
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
