package types

func IsTakeProfitPriceInifite(mtp *MTP) bool {
	return mtp.TakeProfitPrice.TruncateInt().String() == TakeProfitPriceDefault
}
