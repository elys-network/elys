package types

func IsTakeProfitPriceInfinite(mtp MTP) bool {
	return mtp.TakeProfitPrice.Equal(TakeProfitPriceDefault)
}
