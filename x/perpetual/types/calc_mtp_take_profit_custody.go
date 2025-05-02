package types

import (
	"cosmossdk.io/math"
)

func CalcMTPTakeProfitCustody(mtp MTP) math.Int {
	if IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return math.ZeroInt()
	}
	if mtp.Position == Position_LONG {
		return mtp.GetBigDecLiabilities().Quo(mtp.GetBigDecTakeProfitPrice()).Dec().TruncateInt()
	} else {
		return mtp.GetBigDecLiabilities().Mul(mtp.GetBigDecTakeProfitPrice()).Dec().TruncateInt()
	}
}
