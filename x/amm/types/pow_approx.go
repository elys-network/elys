package types

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ln2    = sdk.MustNewDecFromStr("0.693147180559945309")
	invLn2 = sdk.MustNewDecFromStr("1.442695040888963407")
	exp    = sdk.MustNewDecFromStr("2.718281828459045235")
)

func computeExp(x sdk.Dec) (sdk.Dec, error) {
	if x.Equal(sdk.ZeroDec()) {
		return sdk.OneDec(), nil
	}
	if x.Equal(sdk.OneDec()) {
		return exp, nil
	}
	// exp(-42) is approx 5.7 x 10^-19, smallest dec possible is 10^-18
	if x.LTE(sdk.NewDecFromInt(sdk.NewInt(-42))) {
		return sdk.ZeroDec(), nil
	}

	// Range reduction: x = k * ln(2) + y
	k := x.Mul(invLn2).TruncateInt64()
	y := x.Sub(sdk.NewDecFromInt(sdk.NewInt(k)).Mul(ln2))

	expY := sdk.OneDec()
	term := sdk.OneDec()

	//n decides the precision of the value, higher the n, greater is the accuracy
	for n := int64(1); n <= 100; n++ {
		term = term.Mul(y).QuoInt64(n)
		expY = expY.Add(term)
		if term.Abs().LTE(powPrecision) {
			break
		}
	}

	twoPowK := sdk.OneDec()
	if k > 0 {
		twoPowK = two.Power(uint64(k))
	} else if k < 0 {
		twoPowK = sdk.OneDec().Quo(two.Power(uint64(-k)))
	}

	result := expY.Mul(twoPowK)

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
