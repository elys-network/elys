package types

import (
	"cosmossdk.io/math"
)

// ReachedTakeProfitPrice tells if the take profit price is reached
func ReachedTakeProfitPrice(mtp *MTP, assetPrice math.LegacyDec) bool {
	if mtp.Position == Position_LONG {
		return mtp.TakeProfitPrice.GTE(assetPrice)
	} else if mtp.Position == Position_SHORT {
		return mtp.TakeProfitPrice.LTE(assetPrice)
	}
	return false
}
