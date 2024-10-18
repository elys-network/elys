package types

import "cosmossdk.io/math"

func (perpetualPool Pool) GetPerpetualPoolBalancesByPosition(denom string, position Position) (math.Int, math.Int, math.Int, math.Int) {
	poolAsset := perpetualPool.GetPoolAsset(position, denom)
	return poolAsset.Liabilities, poolAsset.Custody, poolAsset.TakeProfitCustody, poolAsset.TakeProfitLiabilities
}

// Get Perpetual Pool Balance
func (perpetualPool Pool) GetPerpetualPoolBalances(denom string) (math.Int, math.Int, math.Int, math.Int) {
	liabilitiesLong, custodyLong, longTakeProfitCustody, longTakeProfitLiabilities := perpetualPool.GetPerpetualPoolBalancesByPosition(denom, Position_LONG)
	liabilitiesShort, custodyShort, shortTakeProfitCustody, shortTakeProfitLiabilities := perpetualPool.GetPerpetualPoolBalancesByPosition(denom, Position_SHORT)

	totalLiabilities := liabilitiesLong.Add(liabilitiesShort)
	totalCustody := custodyLong.Add(custodyShort)
	totalTakeProfitCustody := longTakeProfitCustody.Add(shortTakeProfitCustody)
	totalTakeProfitLiabilities := longTakeProfitLiabilities.Add(shortTakeProfitLiabilities)

	return totalLiabilities, totalCustody, totalTakeProfitCustody, totalTakeProfitLiabilities
}
