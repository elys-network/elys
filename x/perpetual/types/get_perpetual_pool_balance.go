package types

import "cosmossdk.io/math"

func GetPerpetualPoolBalancesByPosition(perpetualPool Pool, denom string, position Position) (math.Int, math.Int, math.Int) {
	poolAsset := perpetualPool.GetPoolAsset(position, denom)
	return poolAsset.AssetBalance, poolAsset.Liabilities, poolAsset.Custody
}

// Get Perpetual Pool Balance
func GetPerpetualPoolBalances(perpetualPool Pool, denom string) (math.Int, math.Int, math.Int) {
	assetBalanceLong, liabilitiesLong, custodyLong := GetPerpetualPoolBalancesByPosition(perpetualPool, denom, Position_LONG)
	assetBalanceShort, liabilitiesShort, custodyShort := GetPerpetualPoolBalancesByPosition(perpetualPool, denom, Position_SHORT)

	assetBalance := assetBalanceLong.Add(assetBalanceShort)
	liabilities := liabilitiesLong.Add(liabilitiesShort)
	custody := custodyLong.Add(custodyShort)

	return assetBalance, liabilities, custody
}
