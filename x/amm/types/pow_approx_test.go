package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
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
		base           sdk.Dec
		exp            sdk.Dec
		expectedResult sdk.Dec
	}{
		{
			// medium base, small exp
			base:           sdk.MustNewDecFromStr("0.8"),
			exp:            sdk.MustNewDecFromStr("0.32"),
			expectedResult: sdk.MustNewDecFromStr("0.93108385"),
		},
		{
			// zero exp
			base:           sdk.MustNewDecFromStr("0.8"),
			exp:            sdk.ZeroDec(),
			expectedResult: sdk.OneDec(),
		},
		{
			// zero base, this should panic
			base:        sdk.ZeroDec(),
			exp:         sdk.OneDec(),
			expectPanic: false,
			expectErr:   true,
		},
		{
			// large base, small exp
			base:           sdk.MustNewDecFromStr("1.9999"),
			exp:            sdk.MustNewDecFromStr("0.23"),
			expectedResult: sdk.MustNewDecFromStr("1.172821461"),
		},
		{
			// large base, large integer exp
			base:           sdk.MustNewDecFromStr("1.777"),
			exp:            sdk.MustNewDecFromStr("20"),
			expectedResult: sdk.MustNewDecFromStr("98570.862372081602"),
		},
		{
			// medium base, large exp, high precision
			base:           sdk.MustNewDecFromStr("1.556"),
			exp:            sdk.MustNewDecFromStr("0.9123"),
			expectedResult: sdk.MustNewDecFromStr("1.4968226674708064"),
		},
		{
			// high base, large exp, high precision
			base:           sdk.MustNewDecFromStr("1.886"),
			exp:            sdk.MustNewDecFromStr("1.9123"),
			expectedResult: sdk.MustNewDecFromStr("3.364483251631"),
		},
		{
			// base equal one
			base:           sdk.MustNewDecFromStr("1"),
			exp:            sdk.MustNewDecFromStr("123"),
			expectedResult: sdk.OneDec(),
		},
		{
			// base close to 2

			base: sdk.MustNewDecFromStr("1.999999999999999999"),
			exp:  sdk.SmallestDec(),
			// In Python: 1.000000000000000000693147181
			expectedResult: sdk.OneDec(),
		},
		{
			// base close to 2 and hitting iteration bound

			base: sdk.MustNewDecFromStr("1.999999999999999999"),
			exp:  sdk.MustNewDecFromStr("0.1"),

			// In Python: 1.071773462536293164

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base: sdk.MustNewDecFromStr("1.99999"),
			exp:  sdk.MustNewDecFromStr("0.1"),

			// expectedResult: sdk.MustNewDecFromStr("1.071772926648356147"),

			// In Python: 1.071772926648356147102864087

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base: sdk.MustNewDecFromStr("1.9999"),
			exp:  sdk.MustNewDecFromStr("0.1"),

			// In Python: 1.071768103548402149880477100
			expectedResult: sdk.MustNewDecFromStr("1.071768103548402149"),
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
