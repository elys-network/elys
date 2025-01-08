package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestDec34(t *testing.T) {
	// Test constructors
	require.Equal(t, NewDec34FromInt64(1), NewDec34FromInt64(1))
	require.Equal(t, NewDec34FromInt64(1), OneDec34())
	require.Equal(t, NewDec34FromInt64(2), TwoDec34())
	require.Equal(t, NewDec34FromInt64(3), ThreeDec34())
	require.Equal(t, NewDec34FromInt64(4), FourDec34())
	require.Equal(t, NewDec34FromInt64(5), FiveDec34())
	require.Equal(t, NewDec34FromInt64(-1), MinusOneDec34())
	require.Equal(t, NewDec34FromInt64(0), ZeroDec34())

	// Test NewDec34FromLegacyDec
	legacyDec := math.LegacyNewDec(123)
	dec34FromLegacy := NewDec34FromLegacyDec(legacyDec)
	require.Equal(t, "123.000000000000000000", dec34FromLegacy.String())

	// Test NewDec34FromInt
	intVal := math.NewInt(456)
	dec34FromInt := NewDec34FromInt(intVal)
	require.Equal(t, "456", dec34FromInt.String())

	// Test arithmetic operations
	one := OneDec34()
	two := TwoDec34()
	three := ThreeDec34()

	// Test Add
	sum := one.Add(two)
	require.Equal(t, three, sum)

	// Test AddInt
	require.Equal(t, NewDec34FromInt64(3), one.AddInt(math.NewInt(2)))

	// Test AddInt64
	require.Equal(t, NewDec34FromInt64(3), one.AddInt64(2))

	// Test AddLegacyDec
	require.Equal(t, NewDec34FromLegacyDec(math.LegacyNewDec(3)), one.AddLegacyDec(math.LegacyNewDec(2)))

	// Test Sub
	diff := three.Sub(two)
	require.Equal(t, one, diff)

	// Test SubInt
	require.Equal(t, NewDec34FromInt64(1), three.SubInt(math.NewInt(2)))

	// Test SubInt64
	require.Equal(t, NewDec34FromInt64(1), three.SubInt64(2))

	// Test SubLegacyDec
	require.Equal(t, NewDec34FromLegacyDec(math.LegacyNewDec(1)), three.SubLegacyDec(math.LegacyNewDec(2)))

	// Test Mul
	prod := two.Mul(three)
	require.Equal(t, NewDec34FromInt64(6), prod)

	// Test MulInt
	require.Equal(t, NewDec34FromInt64(10000).String(), NewDec34WithPrec(1, 2).MulInt(math.NewInt(100)).String())

	// Test MulInt64
	require.Equal(t, NewDec34FromInt64(10000).String(), NewDec34WithPrec(1, 2).MulInt64(100).String())

	// Test MulLegacyDec
	require.Equal(t, NewDec34FromString("10000.0000000000000000").String(), NewDec34WithPrec(1, 2).MulLegacyDec(math.LegacyNewDec(100)).String())

	// Test Quo
	quot := three.Quo(two)
	require.Equal(t, "1.5", quot.String())

	// Test QuoInt
	require.Equal(t, NewDec34FromInt64(1).String(), NewDec34WithPrec(1, 2).QuoInt(math.NewInt(100)).String())

	// Test QuoInt64
	require.Equal(t, NewDec34FromInt64(1).String(), NewDec34WithPrec(1, 2).QuoInt64(100).String())

	// Test QuoLegacyDec
	require.Equal(t, OneDec34().String(), NewDec34WithPrec(1, 2).QuoLegacyDec(math.LegacyNewDec(100)).String())

	// Test division by zero panic
	require.Panics(t, func() {
		three.Quo(ZeroDec34())
	})

	// Test ToLegacyDec
	legacyResult := three.ToLegacyDec()
	require.Equal(t, math.LegacyNewDec(3), legacyResult)

	// Test ToInt
	intResult := three.ToInt()
	require.Equal(t, math.NewInt(3), intResult)

	// Test String
	require.Equal(t, "3", three.String())

	// Test Int64
	require.Equal(t, int64(3), three.ToInt64())

	// Test Equal
	require.True(t, three.Equal(three))
	require.False(t, three.Equal(two))

	// Test IsZero
	require.False(t, three.IsZero())
	require.True(t, ZeroDec34().IsZero())

	// Test IsNegative
	require.False(t, three.IsNegative())
	require.True(t, MinusOneDec34().IsNegative())

	// Test IsPositive
	require.True(t, three.IsPositive())
	require.False(t, MinusOneDec34().IsPositive())

	// Test NewDec34WithPrec
	require.Equal(t, NewDec34WithPrec(1, 2).String(), NewDec34WithPrec(1, 2).String())
	require.Equal(t, NewDec34FromInt64(100).String(), NewDec34WithPrec(1, 2).Mul(NewDec34WithPrec(1, 0)).String())

	// Test NewDec34FromString
	require.Equal(t, NewDec34FromString("1.234567890123456789").String(), NewDec34FromString("1.234567890123456789").String())
	require.Equal(t,
		NewDec34FromString("1000000000000000000000000000000000000123456789.00000000").String(),
		NewDec34FromString("1000000000000000000000000000000000000123456789.00000000").String(),
	)

	// Test Abs
	require.Equal(t, NewDec34FromInt64(1).String(), NewDec34FromInt64(-1).Abs().String())
	require.Equal(t, NewDec34FromInt64(0).String(), ZeroDec34().Abs().String())
	require.Equal(t,
		NewDec34FromString("1000000000000000000000000000000000000123456789.00000000").String(),
		NewDec34FromString("-1000000000000000000000000000000000000123456789.00000000").Abs().String(),
	)

	// Test Neg
	require.Equal(t, NewDec34FromInt64(-1).String(), NewDec34FromInt64(1).Neg().String())
	require.Equal(t, NewDec34FromInt64(0).String(), ZeroDec34().Neg().String())
	require.Equal(t,
		NewDec34FromString("-1000000000000000000000000000000000000123456789.00000000").String(),
		NewDec34FromString("1000000000000000000000000000000000000123456789.00000000").Neg().String(),
	)

	// Test GT
	require.True(t, NewDec34FromInt64(2).GT(NewDec34FromInt64(1)))
	require.False(t, NewDec34FromInt64(1).GT(NewDec34FromInt64(2)))

	// Test LT
	require.True(t, NewDec34FromInt64(1).LT(NewDec34FromInt64(2)))
	require.False(t, NewDec34FromInt64(2).LT(NewDec34FromInt64(1)))

	// Test GTE
	require.True(t, NewDec34FromInt64(2).GTE(NewDec34FromInt64(1)))
	require.False(t, NewDec34FromInt64(1).GTE(NewDec34FromInt64(2)))
	require.True(t, NewDec34FromInt64(2).GTE(NewDec34FromInt64(2)))

	// Test LTE
	require.True(t, NewDec34FromInt64(1).LTE(NewDec34FromInt64(2)))
	require.False(t, NewDec34FromInt64(2).LTE(NewDec34FromInt64(1)))
	require.True(t, NewDec34FromInt64(2).LTE(NewDec34FromInt64(2)))

	// Test MinDec34
	require.Equal(t, NewDec34FromInt64(1), MinDec34(NewDec34FromInt64(1), NewDec34FromInt64(2)))
	require.Equal(t, NewDec34FromInt64(1), MinDec34(NewDec34FromInt64(2), NewDec34FromInt64(1)))

	// Test MaxDec34
	require.Equal(t, NewDec34FromInt64(2), MaxDec34(NewDec34FromInt64(1), NewDec34FromInt64(2)))
	require.Equal(t, NewDec34FromInt64(2), MaxDec34(NewDec34FromInt64(2), NewDec34FromInt64(1)))
}
