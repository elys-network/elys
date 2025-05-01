package utils

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"testing"

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
		base           sdkmath.LegacyDec
		exp            sdkmath.LegacyDec
		expectedResult sdkmath.LegacyDec
	}{
		{
			// medium base, small exp
			base:           sdkmath.LegacyMustNewDecFromStr("0.8"),
			exp:            sdkmath.LegacyMustNewDecFromStr("0.32"),
			expectedResult: sdkmath.LegacyMustNewDecFromStr("0.93108385"),
		},
		{
			// zero exp
			base:           sdkmath.LegacyMustNewDecFromStr("0.8"),
			exp:            sdkmath.LegacyZeroDec(),
			expectedResult: sdkmath.LegacyOneDec(),
		},
		{
			// zero base, this should panic
			base:        sdkmath.LegacyZeroDec(),
			exp:         sdkmath.LegacyOneDec(),
			expectPanic: false,
			expectErr:   true,
		},
		{
			// large base, small exp
			base:           sdkmath.LegacyMustNewDecFromStr("1.9999"),
			exp:            sdkmath.LegacyMustNewDecFromStr("0.23"),
			expectedResult: sdkmath.LegacyMustNewDecFromStr("1.172821461"),
		},
		{
			// large base, large integer exp
			base:           sdkmath.LegacyMustNewDecFromStr("1.777"),
			exp:            sdkmath.LegacyMustNewDecFromStr("20"),
			expectedResult: sdkmath.LegacyMustNewDecFromStr("98570.862372081602"),
		},
		{
			// medium base, large exp, high precision
			base:           sdkmath.LegacyMustNewDecFromStr("1.556"),
			exp:            sdkmath.LegacyMustNewDecFromStr("0.9123"),
			expectedResult: sdkmath.LegacyMustNewDecFromStr("1.4968226674708064"),
		},
		{
			// high base, large exp, high precision
			base:           sdkmath.LegacyMustNewDecFromStr("1.886"),
			exp:            sdkmath.LegacyMustNewDecFromStr("1.9123"),
			expectedResult: sdkmath.LegacyMustNewDecFromStr("3.364483251631"),
		},
		{
			// base equal one
			base:           sdkmath.LegacyMustNewDecFromStr("1"),
			exp:            sdkmath.LegacyMustNewDecFromStr("123"),
			expectedResult: sdkmath.LegacyOneDec(),
		},
		{
			// base close to 2

			base: sdkmath.LegacyMustNewDecFromStr("1.999999999999999999"),
			exp:  sdkmath.LegacySmallestDec(),
			// In Python: 1.000000000000000000693147181
			expectedResult: sdkmath.LegacyOneDec(),
		},
		{
			// base close to 2 and hitting iteration bound

			base: sdkmath.LegacyMustNewDecFromStr("1.999999999999999999"),
			exp:  sdkmath.LegacyMustNewDecFromStr("0.1"),

			// In Python: 1.071773462536293164

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base: sdkmath.LegacyMustNewDecFromStr("1.99999"),
			exp:  sdkmath.LegacyMustNewDecFromStr("0.1"),

			// expectedResult: sdkmath.LegacyMustNewDecFromStr("1.071772926648356147"),

			// In Python: 1.071772926648356147102864087

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base: sdkmath.LegacyMustNewDecFromStr("1.9999"),
			exp:  sdkmath.LegacyMustNewDecFromStr("0.1"),

			// In Python: 1.071768103548402149880477100
			expectedResult: sdkmath.LegacyMustNewDecFromStr("1.071768103548402149"),
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
