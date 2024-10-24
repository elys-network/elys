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

func (k Keeper) GetFundingPaymentRates(ctx sdk.Context, pool types.Pool) (long sdk.Dec, short sdk.Dec) {
	fundingRateLong, fundingRateShort := k.ComputeFundingRate(ctx, pool)

	// account custody from long position
	totalCustodyLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		totalCustodyLong = totalCustodyLong.Add(asset.Custody)
	}

	// account custody from short position
	totalLiabilitiesShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		totalLiabilitiesShort = totalLiabilitiesShort.Add(asset.Liabilities)
	}

	if fundingRateLong.IsZero() {
		// short will pay
		// long will receive
		unpopular_rate := fundingRateShort.Mul(totalLiabilitiesShort.ToLegacyDec()).Quo(totalCustodyLong.ToLegacyDec())
		return unpopular_rate, fundingRateShort.Neg()
	} else {
		// long will pay
		// short will receive
		unpopular_rate := fundingRateLong.Mul(totalCustodyLong.ToLegacyDec()).Quo(totalLiabilitiesShort.ToLegacyDec())
		return fundingRateLong.Neg(), unpopular_rate
	}
}
