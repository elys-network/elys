package types

import "github.com/osmosis-labs/osmosis/osmomath"

// ReachedTakeProfitPrice tells if the take profit price is reached
func ReachedTakeProfitPrice(mtp *MTP, assetPrice osmomath.BigDec) bool {
	if mtp.Position == Position_LONG {
		return mtp.GetBigDecTakeProfitPrice().GTE(assetPrice)
	} else if mtp.Position == Position_SHORT {
		return mtp.GetBigDecTakeProfitPrice().LTE(assetPrice)
	}
	return false
}
