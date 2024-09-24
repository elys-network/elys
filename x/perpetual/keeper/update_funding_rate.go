package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// UpdateFundingRate updates the funding rate of a pool
func (k Keeper) UpdateFundingRate(ctx sdk.Context, pool *types.Pool) error {
	poolAssetsLong := pool.GetPoolAssets(types.Position_LONG)
	poolAssetsShort := pool.GetPoolAssets(types.Position_SHORT)

	uusdc, found := k.assetProfileKeeper.GetEntry(ctx, "uusdc")
	if !found {
		return types.ErrDenomNotFound
	}

	liabilitiesLong := sdk.ZeroInt()
	for _, asset := range *poolAssetsLong {

		if uusdc.Denom == asset.AssetDenom {
			liabilitiesLong = liabilitiesLong.Add(asset.Liabilities)
		} else {
			if asset.Liabilities.IsZero() {
				continue
			}
			coin := sdk.NewCoin(asset.AssetDenom, asset.Liabilities)
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, asset.AssetDenom)

			if err != nil {
				return err
			}
			l, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, coin, uusdc.Denom, ammPool)
			if err != nil {
				return err
			}
			liabilitiesLong = liabilitiesLong.Add(l)
		}

	}

	liabilitiesShort := sdk.ZeroInt()
	for _, asset := range *poolAssetsShort {
		if uusdc.Denom == asset.AssetDenom {
			liabilitiesShort = liabilitiesShort.Add(asset.Liabilities)
		} else {
			if asset.Liabilities.IsZero() {
				continue
			}
			coin := sdk.NewCoin(asset.AssetDenom, asset.Liabilities)
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, asset.AssetDenom)

			if err != nil {
				return err
			}
			l, err := k.OpenShortChecker.EstimateSwapGivenOut(ctx, coin, uusdc.Denom, ammPool)
			if err != nil {
				return err
			}
			liabilitiesShort = liabilitiesShort.Add(l)
		}
	}

	// get params
	params := k.GetParams(ctx)

	// calculate and update funding fee
	pool.FundingRate = types.CalcFundingRate(liabilitiesLong, liabilitiesShort, params.FundingFeeBaseRate, params.FundingFeeMaxRate, params.FundingFeeMinRate)

	return nil
}
