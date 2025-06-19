package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) GetFundingPaymentRates(ctx sdk.Context, pool types.Pool) (long math.LegacyDec, short math.LegacyDec) {
	fundingRateLong, fundingRateShort := k.ComputeFundingRate(ctx, pool)

	totalLongOpenInterest := pool.GetTotalLongOpenInterest()
	totalShortOpenInterest := pool.GetTotalShortOpenInterest()

	if fundingRateLong.IsZero() {
		// short will pay
		// long will receive
		unpopular_rate := math.LegacyZeroDec()
		if !totalLongOpenInterest.IsZero() {
			unpopular_rate = fundingRateShort.Mul(totalShortOpenInterest.ToLegacyDec()).Quo(totalLongOpenInterest.ToLegacyDec())
		}
		return unpopular_rate.Neg(), fundingRateShort
	} else {
		// long will pay
		// short will receive
		unpopular_rate := math.LegacyZeroDec()
		if !totalShortOpenInterest.IsZero() {
			unpopular_rate = fundingRateLong.Mul(totalLongOpenInterest.ToLegacyDec()).Quo(totalShortOpenInterest.ToLegacyDec())
		}
		return fundingRateLong, unpopular_rate.Neg()
	}
}
