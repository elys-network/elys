package types

func GetTradingAsset(collateralAsset string, borrowAsset string, baseCurrency string) string {
	if collateralAsset == baseCurrency {
		return borrowAsset
	}
	return collateralAsset
}
