package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPerpetualPoolBalancesByPosition(perpetualPool Pool, denom string, position Position) (sdk.Int, sdk.Int, sdk.Int) {
	poolAsset := perpetualPool.GetPoolAsset(position, denom)
	return poolAsset.AssetBalance, poolAsset.Liabilities, poolAsset.Custody
}

// Get Perpetual Pool Balance
func GetPerpetualPoolBalances(perpetualPool Pool, denom string) (sdk.Int, sdk.Int, sdk.Int) {
	assetBalanceLong, liabilitiesLong, custodyLong := GetPerpetualPoolBalancesByPosition(perpetualPool, denom, Position_LONG)
	assetBalanceShort, liabilitiesShort, custodyShort := GetPerpetualPoolBalancesByPosition(perpetualPool, denom, Position_SHORT)

	assetBalance := assetBalanceLong.Add(assetBalanceShort)
	liabilities := liabilitiesLong.Add(liabilitiesShort)
	custody := custodyLong.Add(custodyShort)

	return assetBalance, liabilities, custody
}
