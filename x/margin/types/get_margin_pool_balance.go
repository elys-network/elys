package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetMarginPoolBalancesByPosition(marginPool Pool, denom string, position Position) (sdk.Int, sdk.Int, sdk.Int) {
	poolAssets := marginPool.GetPoolAssets(position)

	for _, asset := range *poolAssets {
		if asset.AssetDenom == denom {
			return asset.AssetBalance, asset.Liabilities, asset.Custody
		}
	}

	return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
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
