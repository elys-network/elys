package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// GetNetOpenInterest calculates the net open interest for a given pool.
// Note: Net open interest should always be in terms of trading asset
func (k Keeper) GetNetOpenInterest(ctx sdk.Context, pool types.Pool) math.Int {
	// account custody from long position
	totalCustodyLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		totalCustodyLong = totalCustodyLong.Add(asset.Custody)
	}

	// account liabilities from short position
	totalLiabilityShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		totalLiabilityShort = totalLiabilityShort.Add(asset.Liabilities)
	}

	// Net Open Interest = Long custody - Short Liabilities
	netOpenInterest := totalCustodyLong.Sub(totalLiabilityShort)

	return netOpenInterest
}
