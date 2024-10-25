package types

import (
	"cosmossdk.io/math"
)

func CalcMTPTakeProfitCustody(mtp MTP) math.Int {
	if IsTakeProfitPriceInfinite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return math.ZeroInt()
	}
	if mtp.Position == Position_LONG {
		return mtp.Liabilities.ToLegacyDec().Quo(mtp.TakeProfitPrice).TruncateInt()
	} else {
		return mtp.Liabilities.ToLegacyDec().Mul(mtp.TakeProfitPrice).TruncateInt()
	}
}
