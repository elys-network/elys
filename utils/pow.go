package utils

import (
	sdkmath "cosmossdk.io/math"
	"fmt"

	"github.com/osmosis-labs/osmosis/osmomath"
)

var (
	oneHalf           = sdkmath.LegacyMustNewDecFromStr("0.5")
	twoDec            = sdkmath.LegacyMustNewDecFromStr("2")
	ln2               = sdkmath.LegacyMustNewDecFromStr("0.693147180559945309")
	inverseLn2        = sdkmath.LegacyMustNewDecFromStr("1.442695040888963407")
	euler             = sdkmath.LegacyMustNewDecFromStr("2.718281828459045235")
	powIterationLimit = int64(150_000)

	// PowPrecision Don't EVER change after initializing
	// TODO: Analyze choice here.
	powPrecision = sdkmath.LegacyMustNewDecFromStr("0.00000001")
)

// Pow computes base^(exp)
// However since the exponent is not an integer, we must do an approximation algorithm.
// TODO: In the future, lets add some optimized routines for common exponents, e.g. for common wIn / wOut ratios
// Many simple exponents like 2:1 pools.
func Pow(base osmomath.BigDec, exp osmomath.BigDec) osmomath.BigDec {
	// Exponentiation of a negative base with an arbitrary real exponent is not closed within the reals.
	// You can see this by recalling that `i = (-1)^(.5)`. We have to go to complex numbers to define this.
	// (And would have to implement complex logarithms)
	// We don't have a need for negative bases, so we don't include any such logic.
	if !base.IsPositive() {
		panic(fmt.Errorf("base must be greater than 0, base: %s", base.String()))
	}

	// We will use an approximation algorithm to compute the power.
	// Since computing an integer power is easy, we split up the exponent into
	// an integer component and a fractional component.
	integer := exp.TruncateDec()
	fractional := exp.Sub(integer)

	integerPow := base.PowerInteger(uint64(integer.TruncateInt64()))

	if fractional.IsZero() {
		return integerPow
	}

	fractionalPow, err := powerApproximation(base.Dec(), fractional.Dec())
	if err != nil {
		panic(err)
	}

	return integerPow.Mul(osmomath.BigDecFromDec(fractionalPow))
}
