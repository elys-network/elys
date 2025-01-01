package types

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"fmt"
	"math"
	"strconv"
)

func computeExp(x sdkmath.LegacyDec) (sdkmath.LegacyDec, error) {
	if x.Equal(sdkmath.LegacyZeroDec()) {
		return sdkmath.LegacyOneDec(), nil
	}
	if x.Equal(sdkmath.LegacyOneDec()) {
		return euler, nil
	}
	// exp(-42) is approx 5.7 x 10^-19, smallest dec possible is 10^-18
	if x.LTE(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(-42))) {
		return sdkmath.LegacyZeroDec(), nil
	}

	// Range reduction: x = k * ln(2) + y
	k := x.Mul(inverseLn2).TruncateInt64()
	y := x.Sub(sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(k)).Mul(ln2))

	expY := sdkmath.LegacyOneDec()
	term := sdkmath.LegacyOneDec()

	//n decides the precision of the value, higher the n, greater is the accuracy
	for n := int64(1); ; n++ {
		term = term.Mul(y).QuoInt64(n)
		expY = expY.Add(term)
		if term.Abs().LTE(powPrecision) {
			break
		}
		if n > powIterationLimit {
			return sdkmath.LegacyDec{}, fmt.Errorf("failed to reach precision within %d iterations while comuting Exp for: %s", powIterationLimit, x.String())
		}
	}

	twoPowK := sdkmath.LegacyOneDec()
	if k > 0 {
		twoPowK = twoDec.Power(uint64(k))
	} else if k < 0 {
		twoPowK = sdkmath.LegacyOneDec().Quo(twoDec.Power(uint64(-k)))
	}

	result := expY.Mul(twoPowK)

	return result, nil
}

func computeLn(x sdkmath.LegacyDec) (result sdkmath.LegacyDec, err error) {
	if x.LTE(sdkmath.LegacyZeroDec()) {
		return sdkmath.LegacyDec{}, errors.New("x for computing it's Ln must be greater than 0")
	}
	if x.Equal(sdkmath.LegacyOneDec()) {
		return sdkmath.LegacyZeroDec(), nil
	}
	if x.Equal(twoDec) {
		return ln2, nil
	}

	// To bring x is in the range [0.5, 2]
	// we use ln(x) = k * ln(2) + ln(z), where z is in [0.5, 2]
	k := 0
	for x.GT(twoDec) {
		x = x.Quo(twoDec)
		k++
	}
	for x.LT(oneHalf) {
		x = x.MulInt64(2)
		k--
	}
	y := x.Sub(sdkmath.LegacyOneDec())
	result = sdkmath.LegacyZeroDec()
	yPower := y

	// maximum value of y is 1
	// Though n doesn't have upper cap, this iteration will break as |y| < 1,
	// if y is very close to 1 it will large number of iterations
	for n := int64(1); ; n++ {
		sign := sdkmath.NewInt(-1)
		if (n+1)%2 == 0 {
			sign = sdkmath.OneInt()
		}
		term := yPower.MulInt(sign).QuoInt64(n)
		result = result.Add(term)
		// This won't work if y > 1 because absolute value of term is (y^n)/n,
		// if y > 1, it's an increasing value
		if powPrecision.GT(term.Abs()) {
			break
		}

		if n > powIterationLimit {
			return sdkmath.LegacyDec{}, fmt.Errorf("failed to reach precision within %d iterations while comuting Ln for: %s", powIterationLimit, x.String())
		}

		yPower = yPower.Mul(y)
	}

	return result.Add(ln2.MulInt64(int64(k))), nil
}

// powerApproximation Check exponentialLogarithmicMethod and maclaurinSeriesApproximation to understand the limits of this function
func powerApproximation(base sdkmath.LegacyDec, exp sdkmath.LegacyDec) (sdkmath.LegacyDec, error) {
	if !base.IsPositive() {
		return sdkmath.LegacyDec{}, errors.New("base must be greater than 0")
	}
	if exp.LT(sdkmath.LegacyZeroDec()) {
		return sdkmath.LegacyDec{}, errors.New("exp must be greater than 0")
	}
	if exp.IsZero() {
		return sdkmath.LegacyOneDec(), nil
	}
	if exp.Equal(sdkmath.LegacyOneDec()) {
		return base, nil
	}
	if exp.Equal(sdkmath.LegacyOneDec().Neg()) {
		return sdkmath.LegacyOneDec().Quo(base), nil
	}
	if exp.Equal(oneHalf) {
		output, err := base.ApproxSqrt()
		if err != nil {
			return sdkmath.LegacyDec{}, err
		}
		return output, nil
	}
	// case where exp can be represented as uint64
	if exp.IsInteger() && exp.IsPositive() && exp.LTE(sdkmath.LegacyMustNewDecFromStr(strconv.FormatUint(math.MaxUint64, 10))) {
		return base.Power(uint64(exp.TruncateInt64())), nil
	}

	if exp.GT(sdkmath.LegacyOneDec()) {
		return Pow(base, exp), nil
	}

	if base.GTE(oneHalf) && base.LT(twoDec) {
		return maclaurinSeriesApproximation(base, exp, powPrecision), nil
	}

	return exponentialLogarithmicMethod(base, exp)
}

// exponentialLogarithmicMethod This function can operate on any base value >0,
// but it loses its accuracy when base is close to 2^k where k is an integer
// Error propagation is also an issue, when computeLn and computeExp both calculates upto 8 decimal places,
// if the base is large enough, the error propagates and decreases the inaccuracy.
// For example, when calculating 1000^2.23, it can calculate 1000^0.23 upto required accuracy but when we multiply this result
// to 1000^2, the error propagates to other decimal places
func exponentialLogarithmicMethod(base sdkmath.LegacyDec, exp sdkmath.LegacyDec) (sdkmath.LegacyDec, error) {
	lnBase, err := computeLn(base)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	expResult, err := computeExp(exp.Mul(lnBase))
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	return expResult, nil
}

// maclaurinSeriesApproximation This function is very accurate when 0.5 <= base < 2, over 2 it panics
// When base is extremely close to 2 then this function might panic as it's unable to reach accuracy within desired iterations
// 0 <= exp < 1.
func maclaurinSeriesApproximation(originalBase sdkmath.LegacyDec, exp sdkmath.LegacyDec, precision sdkmath.LegacyDec) sdkmath.LegacyDec {
	if !originalBase.IsPositive() {
		panic(errors.New("base must be greater than 0"))
	}

	if exp.IsZero() {
		return sdkmath.LegacyOneDec()
	}

	// Common case optimization
	// Optimize for it being equal to one-half
	if exp.Equal(oneHalf) {
		output, err := originalBase.ApproxSqrt()
		if err != nil {
			panic(err)
		}
		return output
	}
	// TODO: Make an approx-equal function, and then check if exp * 3 = 1, and do a check accordingly

	// We compute this via taking the maclaurin series of (1 + x)^a
	// where x = base - 1.
	// The maclaurin series of (1 + x)^a = sum_{k=0}^{infty} binom(a, k) x^k
	// Binom(a, k) takes the natural continuation on the first parameter, namely that
	// Binom(a, k) = N/D, where D = k!, and N = a(a-1)(a-2)...(a-k+1)
	// Next we show that the absolute value of each term is less than the last term.
	// Note that the change in term n's value vs term n + 1 is a multiplicative factor of
	// v_n = x(a - n) / (n+1)
	// So if |v_n| < 1, we know that each term has a lesser impact on the result than the last.
	// For our bounds on |x| < 1, |a| < 1,
	// it suffices to see for what n is |v_n| < 1,
	// in the worst parameterization of x = 1, a = -1.
	// v_n = |(-1 + epsilon - n) / (n+1)|
	// So |v_n| is always less than 1, as n ranges over the integers.
	//
	// Note that term_n of the expansion is 1 * prod_{i=0}^{n-1} v_i
	// The error if we stop the expansion at term_n is:
	// error_n = sum_{k=n+1}^{infty} term_k
	// At this point we further restrict a >= 0, so 0 <= a < 1.
	// Now we take the _INCORRECT_ assumption that if term_n < p, then
	// error_n < p.
	// This assumption is obviously wrong.
	// However our usages of this function don't use the full domain.
	// With a > 0, |x| << 1, and p sufficiently low, perhaps this actually is true.

	// TODO: Check with our parameterization
	// TODO: If there's a bug, balancer is also wrong here :thonk:

	base := originalBase.Clone()
	x, xneg := AbsDifferenceWithSign(base, sdkmath.LegacyOneDec())
	term := sdkmath.LegacyOneDec()
	sum := sdkmath.LegacyOneDec()
	negative := false

	a := exp.Clone()
	bigK := sdkmath.LegacyNewDec(0)
	// TODO: Document this computation via taylor expansion
	for i := int64(1); term.GTE(precision); i++ {
		// At each iteration, we need two values, i and i-1.
		// To avoid expensive big.Int allocation, we reuse bigK variable.
		// On this line, bigK == i-1.
		c, cneg := AbsDifferenceWithSign(a, bigK)
		// On this line, bigK == i.
		bigK.SetInt64(i)
		term.MulMut(c).MulMut(x).QuoMut(bigK)

		// a is mutated on absDifferenceWithSign, reset
		a.Set(exp)

		if term.IsZero() {
			break
		}
		if xneg {
			negative = !negative
		}

		if cneg {
			negative = !negative
		}

		if negative {
			sum.SubMut(term)
		} else {
			sum.AddMut(term)
		}

		if i == powIterationLimit {
			panic(fmt.Errorf("failed to reach precision within %d iterations, best guess: %s for %s^%s", powIterationLimit, sum, originalBase, exp))
		}
	}
	return sum
}
