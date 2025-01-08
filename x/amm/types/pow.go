package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	elystypes "github.com/elys-network/elys/types"
)

// Pow computes base^(exp)
// However since the exponent is not an integer, we must do an approximation algorithm.
// TODO: In the future, lets add some optimized routines for common exponents, e.g. for common wIn / wOut ratios
// Many simple exponents like 2:1 pools.
func Pow(base sdkmath.LegacyDec, exp sdkmath.LegacyDec) sdkmath.LegacyDec {
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

	integerPow := base.Power(uint64(integer.TruncateInt64()))

	if fractional.IsZero() {
		return integerPow
	}

	fractionalPow, err := powerApproximation(base, fractional)
	if err != nil {
		panic(err)
	}

	return integerPow.Mul(fractionalPow)
}

func PowDec34(base elystypes.Dec34, exp elystypes.Dec34) elystypes.Dec34 {
	return elystypes.NewDec34FromLegacyDec(Pow(base.ToLegacyDec(), exp.ToLegacyDec()))
}
