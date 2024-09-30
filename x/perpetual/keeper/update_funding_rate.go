package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// UpdateFundingRate updates the funding rate of a pool
func (k Keeper) UpdateFundingRate(ctx sdk.Context, pool *types.Pool) error {
	poolAssetsLong := pool.GetPoolAssets(types.Position_LONG)
	poolAssetsShort := pool.GetPoolAssets(types.Position_SHORT)

	liabilitiesLong := sdkmath.ZeroInt()
	for _, asset := range *poolAssetsLong {
		liabilitiesLong = liabilitiesLong.Add(asset.Liabilities)
	}

	liabilitiesShort := sdkmath.ZeroInt()
	for _, asset := range *poolAssetsShort {
		liabilitiesShort = liabilitiesShort.Add(asset.Liabilities)
	}

	// get params
	params := k.GetParams(ctx)

	// calculate and update funding fee
	pool.FundingRate = types.CalcFundingRate(liabilitiesLong, liabilitiesShort, params.FundingFeeBaseRate, params.FundingFeeMaxRate, params.FundingFeeMinRate)

	return nil
}
