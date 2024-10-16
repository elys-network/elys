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

	var err error

	// Calculate liabilities for long and short assets using the separate helper function
	assetLiabilitiesLong, err := k.CalcTotalLiabilities(ctx, *poolAssetsLong, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return err
	}

	assetLiabilitiesShort, err := k.CalcTotalLiabilities(ctx, *poolAssetsShort, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return err
	}

	// get params
	params := k.GetParams(ctx)

	// calculate and update funding fee
	pool.FundingRate = types.CalcFundingRate(assetLiabilitiesLong, assetLiabilitiesShort, params.FundingFeeBaseRate, params.FundingFeeMaxRate, params.FundingFeeMinRate)

	return nil
}
