package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// GetNetOpenInterest calculates the net open interest for a given pool.
func (k Keeper) GetNetOpenInterest(ctx sdk.Context, pool types.Pool) math.Int {
	uusdc, found := k.assetProfileKeeper.GetEntry(ctx, "uusdc")
	if !found {
		return math.ZeroInt()
	}

	var err error

	// Calculate liabilities for long and short assets using the separate helper function
	assetLiabilitiesLong, err := k.CalcTotalLiabilities(ctx, pool.PoolAssetsLong, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return math.ZeroInt()
	}

	assetLiabilitiesShort, err := k.CalcTotalLiabilities(ctx, pool.PoolAssetsShort, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return math.ZeroInt()
	}

	// Net Open Interest = Long Liabilities - Short Liabilities
	netOpenInterest := assetLiabilitiesLong.Sub(assetLiabilitiesShort)

	return netOpenInterest
}
