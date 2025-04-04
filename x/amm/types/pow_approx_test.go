package types

import (
	"fmt"
	"testing"

	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/require"
)

func ConditionalPanic(t *testing.T, expectPanic bool, sut func()) {
	if expectPanic {
		require.Panics(t, sut)
		return
	}
	require.NotPanics(t, sut)
}

func TestPowApprox(t *testing.T) {
	testCases := []struct {
		expectPanic    bool
		expectErr      bool
		base           osmomath.BigDec
		exp            osmomath.BigDec
		expectedResult osmomath.BigDec
	}{
		{
			// medium base, small exp
			base:           osmomath.MustNewBigDecFromStr("0.8"),
			exp:            osmomath.MustNewBigDecFromStr("0.32"),
			expectedResult: osmomath.MustNewBigDecFromStr("0.93108385"),
		},
		{
			// zero exp
			base:           osmomath.MustNewBigDecFromStr("0.8"),
			exp:            osmomath.ZeroBigDec(),
			expectedResult: osmomath.OneBigDec(),
		},
		{
			// zero base, this should panic
			base:        osmomath.ZeroBigDec(),
			exp:         osmomath.OneBigDec(),
			expectPanic: false,
			expectErr:   true,
		},
		{
			// large base, small exp
			base:           osmomath.MustNewBigDecFromStr("1.9999"),
			exp:            osmomath.MustNewBigDecFromStr("0.23"),
			expectedResult: osmomath.MustNewBigDecFromStr("1.172821461"),
		},
		{
			// large base, large integer exp
			base:           osmomath.MustNewBigDecFromStr("1.777"),
			exp:            osmomath.MustNewBigDecFromStr("20"),
			expectedResult: osmomath.MustNewBigDecFromStr("98570.862372081602"),
		},
		{
			// medium base, large exp, high precision
			base:           osmomath.MustNewBigDecFromStr("1.556"),
			exp:            osmomath.MustNewBigDecFromStr("0.9123"),
			expectedResult: osmomath.MustNewBigDecFromStr("1.4968226674708064"),
		},
		{
			// high base, large exp, high precision
			base:           osmomath.MustNewBigDecFromStr("1.886"),
			exp:            osmomath.MustNewBigDecFromStr("1.9123"),
			expectedResult: osmomath.MustNewBigDecFromStr("3.364483251631"),
		},
		{
			// base equal one
			base:           osmomath.MustNewBigDecFromStr("1"),
			exp:            osmomath.MustNewBigDecFromStr("123"),
			expectedResult: osmomath.OneBigDec(),
		},
		{
			// base close to 2

			base: osmomath.MustNewBigDecFromStr("1.999999999999999999"),
			exp:  osmomath.SmallestBigDec(),
			// In Python: 1.000000000000000000693147181
			expectedResult: osmomath.OneBigDec(),
		},
		{
			// base close to 2 and hitting iteration bound

			base: osmomath.MustNewBigDecFromStr("1.999999999999999999"),
			exp:  osmomath.MustNewBigDecFromStr("0.1"),

			// In Python: 1.071773462536293164

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base: osmomath.MustNewBigDecFromStr("1.99999"),
			exp:  osmomath.MustNewBigDecFromStr("0.1"),

			// expectedResult: osmomath.MustNewBigDecFromStr("1.071772926648356147"),

			// In Python: 1.071772926648356147102864087

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base: osmomath.MustNewBigDecFromStr("1.9999"),
			exp:  osmomath.MustNewBigDecFromStr("0.1"),

			// In Python: 1.071768103548402149880477100
			expectedResult: osmomath.MustNewBigDecFromStr("1.071768103548402149"),
		},
	}

	for i, tc := range testCases {
		ConditionalPanic(t, tc.expectPanic, func() {
			fmt.Println(tc.base)
			actualResult, err := powerApproximation(tc.base, tc.exp)
			if !tc.expectErr {
				require.True(
					t,
					tc.expectedResult.Sub(actualResult).Abs().LTE(powPrecision),
					fmt.Sprintf("test %d failed: expected value & actual value's difference should be less than precision, expected: %s, actual: %s, precision: %s", i, tc.expectedResult.String(), actualResult.String(), powPrecision.String()),
				)
			} else {
				require.Error(t, err)
			}
		})
	}
}
