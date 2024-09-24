package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetNetOpenInterest(ctx sdk.Context, pool types.Pool) math.Int {
	assetLiabilitiesLong := sdk.ZeroInt()
	assetLiabilitiesShort := sdk.ZeroInt()

	uusdc, found := k.assetProfileKeeper.GetEntry(ctx, "uusdc")

	if !found {
		return sdk.ZeroInt()
	}

	for _, asset := range pool.PoolAssetsLong {

		if asset.Liabilities.IsZero() {
			continue
		}

		if uusdc.Denom == asset.AssetDenom {
			assetLiabilitiesLong = assetLiabilitiesLong.Add(asset.Liabilities)
		} else {

			coin := sdk.NewCoin(asset.AssetDenom, asset.Liabilities)
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, asset.AssetDenom)

			if err != nil {
				return sdk.ZeroInt()
			}

			l, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, coin, uusdc.Denom, ammPool)
			if err != nil {
				return sdk.ZeroInt()
			}

			assetLiabilitiesLong = assetLiabilitiesLong.Add(l)
		}

	}

	for _, asset := range pool.PoolAssetsShort {

		if asset.Liabilities.IsZero() {
			continue
		}

		if uusdc.Denom == asset.AssetDenom {
			assetLiabilitiesShort = assetLiabilitiesShort.Add(asset.Liabilities)
		} else {

			coin := sdk.NewCoin(asset.AssetDenom, asset.Liabilities)
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, asset.AssetDenom)

			if err != nil {
				return sdk.ZeroInt()
			}

			l, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, coin, uusdc.Denom, ammPool)
			if err != nil {
				return sdk.ZeroInt()
			}

			assetLiabilitiesShort = assetLiabilitiesShort.Add(l)
		}

	}

	netOpenInterest := assetLiabilitiesLong.Sub(assetLiabilitiesShort)
	return netOpenInterest
}
