package types

func IsTakeProfitPriceInifite(mtp *MTP) bool {
	return mtp.TakeProfitPrice.Equal(TakeProfitPriceDefault)
}
