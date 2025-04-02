package utils

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/cockroachdb/apd/v3"
	"github.com/shopspring/decimal"
	"strings"
)

const Precision = uint32(34)

var dec128Context = apd.Context{
	Precision:   Precision,
	MaxExponent: math.MaxExponent,
	MinExponent: math.MinExponent,
	Traps:       apd.DefaultTraps,
}

var (
	ZeroDec  = math.NewDecFromInt64(0)
	OneDec   = math.NewDecFromInt64(1)
	MinusOne = math.NewDecFromInt64(-1)
)

func MinDec(d1 math.Dec, d2 math.Dec) math.Dec {
	// d1 < d2
	if d1.Cmp(d2) == -1 {
		return d1
	}
	return d2
}

func Abs(d1 math.Dec) math.Dec {
	if d1.IsPositive() || d1.IsZero() {
		return d1
	} else {
		d2, err := d1.Mul(math.NewDecFromInt64(-1))
		if err != nil {
			panic(err)
		}
		return d2
	}
}

func Neg(d1 math.Dec) math.Dec {
	d2, err := d1.Mul(math.NewDecFromInt64(-1))
	if err != nil {
		panic(err)
	}
	return d2
}

func GetPaddedDecString(price math.Dec) string {
	b, err := price.BigInt()
	if err != nil {
		panic(err)
	}
	dec := decimal.NewFromBigInt(b, -1*int32(Precision)).StringFixed(int32(Precision))
	return getPaddedPriceFromString(dec)
}

func getPaddedPriceFromString(price string) string {
	components := strings.Split(price, ".")
	naturalPart, decimalPart := components[0], components[1]
	return fmt.Sprintf("%032s.%s", naturalPart, decimalPart)
}

func IntToDec(a math.Int) math.Dec {
	aDec, err := math.DecFromLegacyDec(a.ToLegacyDec())
	if err != nil {
		panic(err)
	}
	return aDec
}
