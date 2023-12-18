package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Don't EVER change after initializing
// TODO: Analyze choice here.
var powPrecision, _ = sdk.NewDecFromStr("0.00000001")

var (
	one_half sdk.Dec = sdk.MustNewDecFromStr("0.5")
	one      sdk.Dec = sdk.OneDec()
	two      sdk.Dec = sdk.MustNewDecFromStr("2")

	// https://www.wolframalpha.com/input?i=2.718281828459045235360287471352662498&assumption=%22ClashPrefs%22+-%3E+%7B%22Math%22%7D
	// nolint: unused
	// eulersNumber = sdk.MustNewDecFromStr("2.718281828459045235360287471352662498")
)

// Pow computes base^(exp)
// However since the exponent is not an integer, we must do an approximation algorithm.
// TODO: In the future, lets add some optimized routines for common exponents, e.g. for common wIn / wOut ratios
// Many simple exponents like 2:1 pools.
func Pow(base sdk.Dec, exp sdk.Dec) sdk.Dec {
	// Exponentiation of a negative base with an arbitrary real exponent is not closed within the reals.
	// You can see this by recalling that `i = (-1)^(.5)`. We have to go to complex numbers to define this.
	// (And would have to implement complex logarithms)
	// We don't have a need for negative bases, so we don't include any such logic.
	if !base.IsPositive() {
		panic(fmt.Errorf("base must be greater than 0"))
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

	fractionalPow := PowApprox(base, fractional, powPrecision)

	return integerPow.Mul(fractionalPow)
}
