package types

import "cosmossdk.io/math"

type PriceTick uint64

func (v PriceTick) ToUint64() uint64 {
	return uint64(v)
}

func (v PriceTick) ToInt64() int64 {
	return int64(v)
}

func (v PriceTick) ToPrice() math.LegacyDec {
	return math.LegacyNewDec(v.ToInt64()).QuoInt64(PriceMultiplier)
}
