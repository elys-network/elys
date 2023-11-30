package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetMarginPoolBalancesByPosition(marginPool Pool, denom string, position Position) (sdk.Int, sdk.Int, sdk.Int) {
	poolAsset := marginPool.GetPoolAsset(position, denom)
	return poolAsset.AssetBalance, poolAsset.Liabilities, poolAsset.Custody
}

// Get Margin Pool Balance
func GetMarginPoolBalances(marginPool Pool, denom string) (sdk.Int, sdk.Int, sdk.Int) {
	assetBalanceLong, liabilitiesLong, custodyLong := GetMarginPoolBalancesByPosition(marginPool, denom, Position_LONG)
	assetBalanceShort, liabilitiesShort, custodyShort := GetMarginPoolBalancesByPosition(marginPool, denom, Position_SHORT)

	assetBalance := assetBalanceLong.Add(assetBalanceShort)
	liabilities := liabilitiesLong.Add(liabilitiesShort)
	custody := custodyLong.Add(custodyShort)

	return assetBalance, liabilities, custody
}
