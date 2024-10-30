package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// CalcTotalLiabilities computes the total liabilities for a list of pool assets.
// It processes each asset, adding its liabilities directly if it's denominated in uusdc,
// or estimating the swap value otherwise.
func (k Keeper) CalcTotalLiabilities(
	ctx sdk.Context,
	assets []types.PoolAsset,
	ammPoolId uint64,
	uusdcDenom string,
) (math.Int, error) {
	totalLiabilities := math.ZeroInt()

	for _, asset := range assets {
		// Skip assets with zero liabilities
		if asset.Liabilities.IsZero() {
			continue
		}

		if asset.AssetDenom == uusdcDenom {
			// Directly add liabilities for uusdc
			totalLiabilities = totalLiabilities.Add(asset.Liabilities)
			continue
		}

		// Estimate swap and add to total liabilities
		coin := sdk.NewCoin(asset.AssetDenom, asset.Liabilities)
		ammPool, err := k.GetAmmPool(ctx, ammPoolId)
		if err != nil {
			return math.ZeroInt(), err
		}

		estimatedSwap, _, err := k.EstimateSwapGivenOut(ctx, coin, uusdcDenom, ammPool)
		if err != nil {
			return math.ZeroInt(), err
		}

		totalLiabilities = totalLiabilities.Add(estimatedSwap)
	}

	return totalLiabilities, nil
}
