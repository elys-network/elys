package types

import (
	"cosmossdk.io/math"
)

func CalcMTPTakeProfitCustody(mtp *MTP) math.Int {
	if IsTakeProfitPriceInifite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return mtp.Custody
	}
	return math.LegacyNewDecFromInt(mtp.Liabilities).Quo(mtp.TakeProfitPrice).TruncateInt()
}
