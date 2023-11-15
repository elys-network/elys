package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) GetMarginPoolBalancesByPosition(marginPool types.Pool, denom string, position types.Position) (sdk.Int, sdk.Int, sdk.Int) {
	poolAssets := marginPool.GetPoolAssets(position)

	for _, asset := range *poolAssets {
		if asset.AssetDenom == denom {
			return asset.AssetBalance, asset.Liabilities, asset.Custody
		}
	}

	return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
}

// Get Margin Pool Balance
func (k Keeper) GetMarginPoolBalances(marginPool types.Pool, denom string) (assetBalance sdk.Int, liabilities sdk.Int, custody sdk.Int) {
	assetBalanceLong, liabilitiesLong, custodyLong := k.GetMarginPoolBalancesByPosition(marginPool, denom, types.Position_LONG)
	assetBalanceShort, liabilitiesShort, custodyShort := k.GetMarginPoolBalancesByPosition(marginPool, denom, types.Position_SHORT)

	assetBalance = assetBalanceLong.Add(assetBalanceShort)
	liabilities = liabilitiesLong.Add(liabilitiesShort)
	custody = custodyLong.Add(custodyShort)

	return assetBalance, liabilities, custody
}
