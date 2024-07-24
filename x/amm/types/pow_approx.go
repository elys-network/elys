package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ln2 = sdk.MustNewDecFromStr("0.693147180559945309")

func computeExp(x sdk.Dec) sdk.Dec {
	result := sdk.OneDec()
	factorial := int64(1)

	// n decide the accuracy, cannot be more than 20 because 21! > MaxInt64
	// For increased accuracy, can use sdk.Int but then that would slow down the computation
	// We can also store the value of the factorials to reduce computations
	for n := int64(1); n <= int64(20); n++ {
		factorial = factorial * n
		term := x.Power(uint64(n)).QuoInt64(factorial)
		result = result.Add(term)
	}

	return result
}

func computeLn(x sdk.Dec) sdk.Dec {
	if x.LTE(sdk.ZeroDec()) {
		panic("x must be greater than 0")
	}
	if x.Equal(sdk.OneDec()) {
		return sdk.ZeroDec()
	}

	// To bring x is in the range [0.5, 2]
	// we use ln(x) = k * ln(2) + ln(y), where y is in [0.5, 2]
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
	result := sdk.ZeroDec()
	// Ideally, this loop will not go till 100 iterations. This precision should be reached much before.
	// Setting 100 terms to have an upper cap on the iterations.
	for n := uint64(1); n <= 100; n++ {
		sign := sdk.NewInt(-1)
		if (n+1)%2 == 0 {
			sign = sdk.OneInt()
		}
		term := y.Power(n).MulInt(sign).QuoInt64(int64(n))
		result = result.Add(term)

		if powPrecision.GT(term.Abs()) {
			break
		}
	}

	return result.Add(ln2.MulInt64(int64(k)))
}

// PowerApproximation This uses formula z = exp(b * ComputeLn(a)) for a^b.
// `terms` increases the accuracy of ComputeLn(a)
func PowerApproximation(base sdk.Dec, exp sdk.Dec) sdk.Dec {
	if !base.IsPositive() {
		panic(fmt.Errorf("base must be greater than 0"))
	}
	if exp.LTE(sdk.ZeroDec()) {
		panic(fmt.Errorf("exp must be greater than 0"))
	}
	if exp.IsZero() {
		return sdk.OneDec()
	}
	if exp.Equal(sdk.OneDec()) {
		return base
	}
	lnBase := computeLn(base)
	expResult := computeExp(exp.Mul(lnBase))
	return expResult
}

// Contract: 0 < base <= 2
// 0 <= exp < 1.
func PowApprox(base sdk.Dec, exp sdk.Dec, precision sdk.Dec) sdk.Dec {
	if !base.IsPositive() {
		panic(fmt.Errorf("base must be greater than 0"))
	}

	if exp.IsZero() {
		return sdk.OneDec()
	}

	// Common case optimization
	// Optimize for it being equal to one-half
	if exp.Equal(one_half) {
		output, err := base.ApproxSqrt()
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
	// TODO: If theres a bug, balancer is also wrong here :thonk:

	base = base.Clone()
	x, xneg := AbsDifferenceWithSign(base, one)
	term := sdk.OneDec()
	sum := sdk.OneDec()
	negative := false

	a := exp.Clone()
	bigK := sdk.NewDec(0)
	// TODO: Document this computation via taylor expansion
	for i := int64(1); term.GTE(precision); i++ {
		// At each iteration, we need two values, i and i-1.
		// To avoid expensive big.Int allocation, we reuse bigK variable.
		// On this line, bigK == i-1.
		c, cneg := AbsDifferenceWithSign(a, bigK)
		// On this line, bigK == i.
		bigK.Set(sdk.NewDec(i)) // TODO: O(n) bigint allocation happens
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
	}
	return sum
}
