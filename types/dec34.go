package types

import (
	"cosmossdk.io/math"
	"github.com/cockroachdb/apd/v2"
	regenmath "github.com/regen-network/regen-ledger/types/v2/math"
)

type Dec34 regenmath.Dec

func NewDec34FromString(s string) Dec34 {
	y, err := regenmath.NewDecFromString(s)
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func NewDec34FromLegacyDec(d math.LegacyDec) Dec34 {
	y, err := regenmath.NewDecFromString(d.String())
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func NewDec34FromInt(i math.Int) Dec34 {
	y, err := regenmath.NewDecFromString(i.String())
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func NewDec34FromInt64(i int64) Dec34 {
	return Dec34(regenmath.NewDecFromInt64(i))
}

func NewDec34WithPrec(i int64, prec int32) Dec34 {
	return Dec34(regenmath.NewDecFinite(i, -prec))
}

func ZeroDec34() Dec34     { return NewDec34FromInt64(0) }
func OneDec34() Dec34      { return NewDec34FromInt64(1) }
func TwoDec34() Dec34      { return NewDec34FromInt64(2) }
func ThreeDec34() Dec34    { return NewDec34FromInt64(3) }
func FourDec34() Dec34     { return NewDec34FromInt64(4) }
func FiveDec34() Dec34     { return NewDec34FromInt64(5) }
func MinusOneDec34() Dec34 { return NewDec34FromInt64(-1) }

func MinDec34(a, b Dec34) Dec34 {
	if a.LT(b) {
		return a
	}
	return b
}

func MaxDec34(a, b Dec34) Dec34 {
	if a.GT(b) {
		return a
	}
	return b
}

func PowDec34(base, exp Dec34) Dec34 {
	b, _, _ := apd.NewFromString(base.String())
	e, _, _ := apd.NewFromString(exp.String())
	d := new(apd.Decimal)

	// add enough precision to handle big decimal values
	c := apd.BaseContext.WithPrecision(100)
	_, err := c.Pow(d, b, e)
	if err != nil {
		panic(err)
	}

	return NewDec34FromString(d.Text('f'))
}

func (d Dec34) Pow(exp Dec34) Dec34 {
	return PowDec34(d, exp)
}

func (d Dec34) PowLegacyDec(exp math.LegacyDec) Dec34 {
	return PowDec34(d, NewDec34FromLegacyDec(exp))
}

func (d Dec34) String() string {
	x, _ := regenmath.Dec(d).Reduce()
	return x.String()
}

func (d Dec34) Mul(other Dec34) Dec34 {
	y, err := regenmath.Dec(d).Mul(regenmath.Dec(other))
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func (d Dec34) MulInt(i math.Int) Dec34 {
	return d.Mul(NewDec34FromInt(i))
}

func (d Dec34) MulInt64(i int64) Dec34 {
	return d.Mul(NewDec34FromInt64(i))
}

func (d Dec34) MulLegacyDec(other math.LegacyDec) Dec34 {
	return d.Mul(NewDec34FromLegacyDec(other))
}

func (d Dec34) Quo(other Dec34) Dec34 {
	y, err := regenmath.Dec(d).Quo(regenmath.Dec(other))
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func (d Dec34) QuoInt(i math.Int) Dec34 {
	y, err := regenmath.Dec(d).Quo(regenmath.Dec(NewDec34FromInt(i)))
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func (d Dec34) QuoInt64(i int64) Dec34 {
	return d.QuoInt(math.NewInt(i))
}

func (d Dec34) QuoLegacyDec(other math.LegacyDec) Dec34 {
	return d.Quo(NewDec34FromLegacyDec(other))
}

func (d Dec34) Add(other Dec34) Dec34 {
	y, err := regenmath.Dec(d).Add(regenmath.Dec(other))
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func (d Dec34) AddInt(i math.Int) Dec34 {
	return d.Add(NewDec34FromInt(i))
}

func (d Dec34) AddInt64(i int64) Dec34 {
	return d.Add(NewDec34FromInt64(i))
}

func (d Dec34) AddLegacyDec(other math.LegacyDec) Dec34 {
	return d.Add(NewDec34FromLegacyDec(other))
}

func (d Dec34) Sub(other Dec34) Dec34 {
	y, err := regenmath.Dec(d).Sub(regenmath.Dec(other))
	if err != nil {
		panic(err)
	}
	return Dec34(y)
}

func (d Dec34) SubInt(i math.Int) Dec34 {
	return d.Sub(NewDec34FromInt(i))
}

func (d Dec34) SubInt64(i int64) Dec34 {
	return d.Sub(NewDec34FromInt64(i))
}

func (d Dec34) SubLegacyDec(other math.LegacyDec) Dec34 {
	return d.Sub(NewDec34FromLegacyDec(other))
}

func (d Dec34) Abs() Dec34 {
	x, _, err := apd.NewFromString(d.String())
	if err != nil {
		panic(err)
	}

	x.Abs(x)

	return NewDec34FromString(x.String())
}

func (d Dec34) Neg() Dec34 {
	x, _, err := apd.NewFromString(d.String())
	if err != nil {
		panic(err)
	}

	x.Neg(x)

	return NewDec34FromString(x.String())
}

func (d Dec34) ToLegacyDec() math.LegacyDec {
	y, _, err := apd.NewFromString(d.String())
	if err != nil {
		panic(err)
	}
	z := new(apd.Decimal)

	// add enough precision to handle big decimal values
	c := apd.BaseContext.WithPrecision(100)
	_, err = c.Quantize(z, y, -18)
	if err != nil {
		panic(err)
	}

	return math.LegacyMustNewDecFromStr(z.Text('f'))
}

func (d Dec34) Ceil() Dec34 {
	y, _, err := apd.NewFromString(d.String())
	if err != nil {
		panic(err)
	}
	z := new(apd.Decimal)

	// add enough precision to handle big decimal values
	c := apd.BaseContext.WithPrecision(100)
	_, err = c.Ceil(z, y)
	if err != nil {
		panic(err)
	}

	return NewDec34FromString(z.Text('f'))
}

func (d Dec34) ToInt() math.Int {
	return regenmath.Dec(d).SdkIntTrim()
}

func (d Dec34) ToInt64() int64 {
	y, err := regenmath.Dec(d).Int64()
	if err != nil {
		panic(err)
	}
	return y
}

func (d Dec34) Equal(other Dec34) bool {
	return regenmath.Dec(d).Equal(regenmath.Dec(other))
}

func (d Dec34) GT(other Dec34) bool {
	return regenmath.Dec(d).Cmp(regenmath.Dec(other)) == 1
}

func (d Dec34) LT(other Dec34) bool {
	return regenmath.Dec(d).Cmp(regenmath.Dec(other)) == -1
}

func (d Dec34) GTE(other Dec34) bool {
	return regenmath.Dec(d).Cmp(regenmath.Dec(other)) >= 0
}

func (d Dec34) LTE(other Dec34) bool {
	return regenmath.Dec(d).Cmp(regenmath.Dec(other)) <= 0
}

func (d Dec34) IsZero() bool {
	return regenmath.Dec(d).IsZero()
}

func (d Dec34) IsNegative() bool {
	return regenmath.Dec(d).IsNegative()
}

func (d Dec34) IsPositive() bool {
	return regenmath.Dec(d).IsPositive()
}
