package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ln2 = sdk.MustNewDecFromStr("0.693147180559945309")

// ComputeExp For values > 22 this function panics to Int overflow because powerTerm have 77 digits
// To extend the range of the function we can reduce the number of iterations at the cost of accuracy
// Since we are only using to compute power approximation for base^exp, it will overflow if Ln(base) > 22,
// which makes base should not be more than approx. 3.58 x 10^9 and exp is very close to 1
func computeExp(x sdk.Dec) (result sdk.Dec, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = errors.New("out of bounds")
			}
		}
	}()
	result = sdk.OneDec()
	factorial := sdk.OneInt()
	powerTerm := x

	//n decides the precision of the value, higher the n, greater is the accuracy
	//it cannot be more than 57 as 58! > 2^256 and 2^256 is the max value sdk.Int can have because cap on bit length
	for n := 1; n <= 57; n++ {
		factorial = factorial.MulRaw(int64(n))
		term := powerTerm.QuoInt(factorial)
		result = result.Add(term)
		if n < 57 {
			powerTerm = powerTerm.Mul(x)
		}
	}

	return result, nil
}

func computeLn(x sdk.Dec) (result sdk.Dec, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = errors.New("out of bounds")
			}
		}
	}()
	if x.LTE(sdk.ZeroDec()) {
		panic("x must be greater than 0")
	}
	if x.Equal(sdk.OneDec()) {
		return sdk.ZeroDec(), nil
	}

	// To bring x is in the range [0.5, 2]
	// we use ln(x) = k * ln(2) + ln(z), where z is in [0.5, 2]
	k := 0
	for x.GT(two) {
		x = x.Quo(two)
		k++
	}
	for x.LT(sdk.OneDec().Quo(two)) {
		x = x.MulInt64(2)
		k--
	}
	y := x.Sub(sdk.OneDec())
	result = sdk.ZeroDec()
	yPower := y

	// maximum value of y is 1
	// Though n doesn't have upper cap, this iteration will break as |y| < 1,
	// if y is very close to 1 it will large number of iterations
	for n := uint64(1); ; n++ {
		sign := sdk.NewInt(-1)
		if (n+1)%2 == 0 {
			sign = sdk.OneInt()
		}
		term := yPower.MulInt(sign).QuoInt64(int64(n))
		result = result.Add(term)
		// This won't work if y > 1 because absolute value of term is (y^n)/n,
		// if y > 1, it's an increasing value
		if powPrecision.GT(term.Abs()) {
			break
		}

		yPower = yPower.Mul(y)
	}

	return result.Add(ln2.MulInt64(int64(k))), nil
}

// PowerApproximation This uses formula z = exp(b * Ln(a)) for a^b.
// Important: This function fails when b * Ln(a) > 22,
// since we are using it to calculate for 0 < b < 1, the upper cap on a is 3.58 x 10^9
// However accuracy decreases fast when 'a' is large
func PowerApproximation(base sdk.Dec, exp sdk.Dec) (sdk.Dec, error) {
	if !base.IsPositive() {
		return sdk.Dec{}, fmt.Errorf("base must be greater than 0")
	}
	if exp.LTE(sdk.ZeroDec()) {
		return sdk.Dec{}, fmt.Errorf("exp must be greater than 0")
	}
	if exp.IsZero() {
		return sdk.OneDec(), nil
	}
	if exp.Equal(sdk.OneDec()) {
		return base, nil
	}
	if exp.Equal(one_half) {
		output, err := base.ApproxSqrt()
		if err != nil {
			return sdk.Dec{}, err
		}
		return output, nil
	}
	lnBase, err := computeLn(base)
	if err != nil {
		return sdk.Dec{}, err
	}
	expResult, err := computeExp(exp.Mul(lnBase))
	if err != nil {
		return sdk.Dec{}, err
	}
	return expResult, nil
}
