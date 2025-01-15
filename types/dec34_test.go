package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cockroachdb/apd/v2"
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
	require.Equal(t, "123", dec34FromLegacy.String())

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
	require.Equal(t, NewDec34FromInt64(10000).String(), OneDec34().MulInt(math.NewInt(10000)).String())

	// Test MulInt64
	require.Equal(t, NewDec34FromInt64(10000).String(), OneDec34().MulInt64(10000).String())

	// Test MulLegacyDec
	require.Equal(t, NewDec34FromString("10000.000000000000000000").String(), OneDec34().MulLegacyDec(math.LegacyNewDec(10000)).String())

	// Test Quo
	require.Equal(t, "1.5", three.Quo(two).String())

	// Test QuoInt
	require.Equal(t, OneDec34().String(), NewDec34FromInt64(100).QuoInt(math.NewInt(100)).String())

	// Test QuoInt64
	require.Equal(t, OneDec34().String(), NewDec34FromInt64(100).QuoInt64(100).String())
	require.Equal(t, NewDec34FromString("0.1585489599188229325215626585489599").String(), NewDec34FromInt64(1000000).Quo(NewDec34FromString("6307200")).String())
	require.Equal(t, NewDec34FromString("0.1585489599188229325215626585489599").String(), NewDec34FromInt64(1000000).QuoInt(math.NewInt(6307200)).String())

	// Test QuoLegacyDec
	require.Equal(t, OneDec34().String(), NewDec34FromInt64(100).QuoLegacyDec(math.LegacyNewDec(100)).String())

	require.Equal(t,
		NewDec34FromString("100000123456.789000000000000000").String(),
		NewDec34FromString("100000123456789.000000000000000000").Quo(NewDec34FromInt64(1000)).String(),
	)

	// Test division by zero panic
	require.Panics(t, func() {
		three.Quo(ZeroDec34())
	})

	// Test ToLegacyDec
	require.Equal(t, math.LegacyNewDec(3), three.ToLegacyDec())
	require.Equal(t,
		NewDec34FromString("1000000000000000000000000000000000000123456789.000000000000000000").ToLegacyDec().String(),
		math.LegacyMustNewDecFromStr("1000000000000000000000000000000000000123456789.000000000000000000").String(),
	)
	require.Equal(t,
		NewDec34FromString("1.000000000000000000000").ToLegacyDec().String(),
		math.LegacyMustNewDecFromStr("1.000000000000000000").String(),
	)
	require.Equal(t,
		NewDec34FromString("123.456000000000000000000").ToLegacyDec().String(),
		math.LegacyMustNewDecFromStr("123.456000000000000000").String(),
	)
	require.Equal(t,
		NewDec34FromString("0.000000000000000000000").ToLegacyDec().String(),
		math.LegacyMustNewDecFromStr("0").String(),
	)
	require.Equal(t,
		NewDec34FromString("96346.39698847304510148894982122764").ToLegacyDec().String(),
		math.LegacyMustNewDecFromStr("96346.396988473045101489").String(),
	)
	require.Equal(t,
		NewDec34FromString("0.0000002").ToLegacyDec().String(),
		math.LegacyMustNewDecFromStr("0.0000002").String(),
	)
	input, _, _ := apd.NewFromString("96346.39698847304510148894982122764")
	output := new(apd.Decimal)
	c := apd.BaseContext.WithPrecision(34)
	c.Quantize(output, input, -18)
	require.Equal(t,
		"96346.396988473045101489",
		output.Text('f'),
	)

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
	require.Equal(t, NewDec34FromInt64(100).String(), OneDec34().Mul(NewDec34WithPrec(100, 0)).String())
	require.Equal(t, NewDec34WithPrec(1, 2).String()+"0000000000000000", math.LegacyNewDecWithPrec(1, 2).String())
	require.Equal(t, NewDec34WithPrec(1, 2).String(), "0.01")

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

	// Test integer powers
	require.Equal(t, "8", PowDec34(NewDec34FromInt64(2), NewDec34FromInt64(3)).String())
	require.Equal(t, "16", PowDec34(NewDec34FromInt64(2), NewDec34FromInt64(4)).String())
	require.Equal(t, "27", PowDec34(NewDec34FromInt64(3), NewDec34FromInt64(3)).String())

	// Test decimal powers
	require.Equal(t, "1.414213562373095048801688724209698078569671875376948073176679737990732478462107038850387534327641573",
		PowDec34(NewDec34FromInt64(2), NewDec34FromString("0.5")).String())
	require.Equal(t, "3.162277660168379331998893544432718533719555139325216826857504852792594438639238221344248108379300295",
		PowDec34(NewDec34FromInt64(10), NewDec34FromString("0.5")).String())

	// Test powers with decimal base
	require.Equal(t, "3.375", PowDec34(NewDec34FromString("1.5"), NewDec34FromInt64(3)).String())

	// Test power of 1 and 0
	require.Equal(t, "1", PowDec34(NewDec34FromInt64(5), NewDec34FromInt64(0)).String())
	require.Equal(t, "5", PowDec34(NewDec34FromInt64(5), NewDec34FromInt64(1)).String())

	// Test Pow method
	require.Equal(t, "8", NewDec34FromInt64(2).Pow(NewDec34FromInt64(3)).String())
	require.Equal(t, "16", NewDec34FromInt64(2).Pow(NewDec34FromInt64(4)).String())
	require.Equal(t, "27", NewDec34FromInt64(3).Pow(NewDec34FromInt64(3)).String())

	// Test PowLegacyDec method
	require.Equal(t, "8", NewDec34FromInt64(2).PowLegacyDec(math.LegacyNewDec(3)).String())
	require.Equal(t, "16", NewDec34FromInt64(2).PowLegacyDec(math.LegacyNewDec(4)).String())
	require.Equal(t, "27", NewDec34FromInt64(3).PowLegacyDec(math.LegacyNewDec(3)).String())

	// Test panic cases for NewDec34FromString
	require.Panics(t, func() {
		NewDec34FromString("invalid.decimal.string")
	})
	require.Panics(t, func() {
		NewDec34FromString("abc")
	})

	// Test panic cases for NewDec34FromLegacyDec
	require.Panics(t, func() {
		invalidLegacyDec := math.LegacyDec{}  // zero value, invalid state
		NewDec34FromLegacyDec(invalidLegacyDec)
	})

	// Test panic cases for NewDec34FromInt
	require.Panics(t, func() {
		invalidInt := math.Int{}  // zero value, invalid state
		NewDec34FromInt(invalidInt)
	})

	// Test Ceil method
	require.Equal(t, NewDec34FromInt64(1).String(), NewDec34FromInt64(1).Ceil().String())
	require.Equal(t, NewDec34FromInt64(2).String(), NewDec34FromString("1.5").Ceil().String())
	require.Equal(t, NewDec34FromInt64(2).String(), NewDec34FromString("1.9").Ceil().String())
	require.Equal(t, NewDec34FromInt64(2).String(), NewDec34FromString("1.999999999999999999").Ceil().String())
	require.Equal(t, NewDec34FromInt64(3).String(), NewDec34FromString("2.000000000000000001").Ceil().String())
}
