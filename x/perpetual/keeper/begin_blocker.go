package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	currentHeight := ctx.BlockHeight()
	pools := k.GetAllPools(ctx)

	for _, pool := range pools {
		rate, err := k.BorrowInterestRateComputation(ctx, pool)
		if err != nil {
			ctx.Logger().Error(err.Error())
			continue
		}
		pool.BorrowInterestRate = rate.Dec()
		pool.LastHeightBorrowInterestRateComputed = currentHeight

		k.SetBorrowRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.InterestBlock{
			InterestRate: rate.Dec(),
			BlockHeight:  ctx.BlockHeight(),
			BlockTime:    ctx.BlockTime().Unix(),
		})

		err = k.UpdatePoolHealth(ctx, &pool)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}

		fundingRateLong, fundingRateShort := k.ComputeFundingRate(ctx, pool)

		pool.FundingRate = fundingRateLong.Dec()
		if fundingRateLong.IsZero() {
			pool.FundingRate = fundingRateShort.Dec().Neg()
		}

		totalLongOpenInterest := pool.GetTotalLongOpenInterest()
		totalShortOpenInterest := pool.GetTotalShortOpenInterest()

		blocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
		fundingAmountLong := osmomath.BigDecFromSDKInt(types.CalcTakeAmount(totalLongOpenInterest, fundingRateLong)).QuoInt64(blocksPerYear)
		fundingAmountShort := osmomath.BigDecFromSDKInt(types.CalcTakeAmount(totalShortOpenInterest, fundingRateShort)).QuoInt64(blocksPerYear)

		fundingShareLong := osmomath.ZeroBigDec()
		if totalShortOpenInterest.IsPositive() {
			fundingShareLong = fundingAmountLong.Quo(osmomath.BigDecFromSDKInt(totalShortOpenInterest))
		}

		fundingShareShort := osmomath.ZeroBigDec()
		if totalLongOpenInterest.IsPositive() {
			fundingShareShort = fundingAmountShort.Quo(osmomath.BigDecFromSDKInt(totalLongOpenInterest))
		}

		k.SetFundingRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.FundingRateBlock{
			BlockHeight:       ctx.BlockHeight(),
			BlockTime:         ctx.BlockTime().Unix(),
			FundingRateLong:   fundingRateLong.Dec(),
			FundingRateShort:  fundingRateShort.Dec(),
			FundingShareShort: fundingShareShort.Dec(),
			FundingShareLong:  fundingShareLong.Dec(),
		})
		k.SetPool(ctx, pool)
	}
}

func (k Keeper) ComputeFundingRate(ctx sdk.Context, pool types.Pool) (osmomath.BigDec, osmomath.BigDec) {
	// Custody amount for long is trading asset -
	// Liability amount for short is trading asset
	// popular_rate = fixed_rate * abs(Custody-Liability) / (Custody+Liability)
	totalLongOpenInterest := pool.GetTotalLongOpenInterest()
	totalShortOpenInterest := pool.GetTotalShortOpenInterest()

	if totalLongOpenInterest.IsZero() || totalShortOpenInterest.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec()
	}

	fixedRate := k.GetParams(ctx).GetBigDecFixedFundingRate()
	if totalLongOpenInterest.GT(totalShortOpenInterest) {
		// long is popular
		// long pays short
		netLongRatio := osmomath.BigDecFromSDKInt(totalLongOpenInterest.Sub(totalShortOpenInterest)).Quo(osmomath.BigDecFromSDKInt(totalLongOpenInterest.Add(totalShortOpenInterest)))
		return netLongRatio.Mul(fixedRate), osmomath.ZeroBigDec()
	} else {
		// short is popular
		// short pays long
		netShortRatio := osmomath.BigDecFromSDKInt(totalShortOpenInterest.Sub(totalLongOpenInterest)).Quo(osmomath.BigDecFromSDKInt(totalLongOpenInterest.Add(totalShortOpenInterest)))
		return osmomath.ZeroBigDec(), netShortRatio.Mul(fixedRate)
	}
}
