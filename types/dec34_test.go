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

	// Test Sub
	diff := three.Sub(two)
	require.Equal(t, one, diff)

	// Test Mul
	prod := two.Mul(three)
	require.Equal(t, NewDec34FromInt64(6), prod)

	// Test Quo
	quot := three.Quo(two)
	require.Equal(t, "1.5", quot.String())

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
}
