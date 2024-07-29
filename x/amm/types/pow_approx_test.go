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
		base           sdk.Dec
		exp            sdk.Dec
		powPrecision   sdk.Dec
		expectedResult sdk.Dec
	}{
		{
			// medium base, small exp
			base:           sdk.MustNewDecFromStr("0.8"),
			exp:            sdk.MustNewDecFromStr("0.32"),
			powPrecision:   sdk.MustNewDecFromStr("0.00000001"),
			expectedResult: sdk.MustNewDecFromStr("0.93108385"),
		},
		{
			// zero exp
			base:           sdk.MustNewDecFromStr("0.8"),
			exp:            sdk.ZeroDec(),
			powPrecision:   sdk.MustNewDecFromStr("0.00001"),
			expectedResult: sdk.OneDec(),
		},
		{
			// zero base, this should panic
			base:         sdk.ZeroDec(),
			exp:          sdk.OneDec(),
			powPrecision: sdk.MustNewDecFromStr("0.00001"),
			expectPanic:  true,
		},
		{
			// large base, small exp
			base:           sdk.MustNewDecFromStr("1.9999"),
			exp:            sdk.MustNewDecFromStr("0.23"),
			powPrecision:   sdk.MustNewDecFromStr("0.000000001"),
			expectedResult: sdk.MustNewDecFromStr("1.172821461"),
		},
		{
			// large base, large integer exp
			base:           sdk.MustNewDecFromStr("1.777"),
			exp:            sdk.MustNewDecFromStr("20"),
			powPrecision:   sdk.MustNewDecFromStr("0.000000000001"),
			expectedResult: sdk.MustNewDecFromStr("98570.862372081602"),
		},
		{
			// medium base, large exp, high precision
			base:           sdk.MustNewDecFromStr("1.556"),
			exp:            sdk.MustNewDecFromStr("20.9123"),
			powPrecision:   sdk.MustNewDecFromStr("0.0000000000000001"),
			expectedResult: sdk.MustNewDecFromStr("10360.058421529811344618"),
		},
		{
			// high base, large exp, high precision
			base:           sdk.MustNewDecFromStr("1.886"),
			exp:            sdk.MustNewDecFromStr("31.9123"),
			powPrecision:   sdk.MustNewDecFromStr("0.00000000000001"),
			expectedResult: sdk.MustNewDecFromStr("621110716.84727942280335811"),
		},
		{
			// base equal one
			base:           sdk.MustNewDecFromStr("1"),
			exp:            sdk.MustNewDecFromStr("123"),
			powPrecision:   sdk.MustNewDecFromStr("0.00000001"),
			expectedResult: sdk.OneDec(),
		},
		{
			// base close to 2

			base:         sdk.MustNewDecFromStr("1.999999999999999999"),
			exp:          sdk.SmallestDec(),
			powPrecision: powPrecision,
			// In Python: 1.000000000000000000693147181
			expectedResult: sdk.OneDec(),
		},
		{
			// base close to 2 and hitting iteration bound

			base:         sdk.MustNewDecFromStr("1.999999999999999999"),
			exp:          sdk.MustNewDecFromStr("0.1"),
			powPrecision: powPrecision,

			// In Python: 1.071773462536293164

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base:         sdk.MustNewDecFromStr("1.99999"),
			exp:          sdk.MustNewDecFromStr("0.1"),
			powPrecision: powPrecision,

			// expectedResult: sdk.MustNewDecFromStr("1.071772926648356147"),

			// In Python: 1.071772926648356147102864087

			expectPanic: true,
		},
		{
			// base close to 2 under iteration limit

			base:         sdk.MustNewDecFromStr("1.9999"),
			exp:          sdk.MustNewDecFromStr("0.1"),
			powPrecision: powPrecision,

			// In Python: 1.071768103548402149880477100
			expectedResult: sdk.MustNewDecFromStr("1.071768103548402149"),
		},
	}

	for i, tc := range testCases {
		var actualResult sdk.Dec
		ConditionalPanic(t, tc.expectPanic, func() {
			fmt.Println(tc.base)
			actualResult = PowApprox(tc.base, tc.exp, tc.powPrecision)
			require.True(
				t,
				tc.expectedResult.Sub(actualResult).Abs().LTE(tc.powPrecision),
				fmt.Sprintf("test %d failed: expected value & actual value's difference should be less than precision, expected: %s, actual: %s, precision: %s", i, tc.expectedResult, actualResult, tc.powPrecision),
			)
		})
	}
}
