package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetFundingPaymentRates(ctx sdk.Context, pool types.Pool) (long osmomath.BigDec, short osmomath.BigDec) {
	fundingRateLong, fundingRateShort := k.ComputeFundingRate(ctx, pool)

	totalLongOpenInterest := pool.GetTotalLongOpenInterest()
	totalShortOpenInterest := pool.GetTotalShortOpenInterest()

	if fundingRateLong.IsZero() {
		// short will pay
		// long will receive
		unpopular_rate := osmomath.ZeroBigDec()
		if !totalLongOpenInterest.IsZero() {
			unpopular_rate = fundingRateShort.Mul(osmomath.BigDecFromSDKInt(totalShortOpenInterest)).Quo(osmomath.BigDecFromSDKInt(totalLongOpenInterest))
		}
		return unpopular_rate.Neg(), fundingRateShort
	} else {
		// long will pay
		// short will receive
		unpopular_rate := osmomath.ZeroBigDec()
		if !totalShortOpenInterest.IsZero() {
			unpopular_rate = fundingRateLong.Mul(osmomath.BigDecFromSDKInt(totalLongOpenInterest)).Quo(osmomath.BigDecFromSDKInt(totalShortOpenInterest))
		}
		return fundingRateLong, unpopular_rate.Neg()
	}
}
