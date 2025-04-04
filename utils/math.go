package utils

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func GetPaddedDecString(price math.LegacyDec) string {
	dec := decimal.NewFromBigInt(price.BigInt(), -18).StringFixed(math.LegacyPrecision)
	return getPaddedPriceFromString(dec)
}

func getPaddedPriceFromString(price string) string {
	components := strings.Split(price, ".")
	naturalPart, decimalPart := components[0], components[1]
	return fmt.Sprintf("%032s.%s", naturalPart, decimalPart)
}
