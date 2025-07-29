package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
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
		pool.BorrowInterestRate = rate
		pool.LastHeightBorrowInterestRateComputed = currentHeight

		k.SetBorrowRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.InterestBlock{
			InterestRate: rate,
			BlockHeight:  ctx.BlockHeight(),
			BlockTime:    ctx.BlockTime().Unix(),
		})

		err = k.UpdatePoolHealth(ctx, &pool)
		if err != nil {
			ctx.Logger().Error(err.Error())
		}

		fundingRateLong, fundingRateShort := k.ComputeFundingRate(ctx, pool)

		pool.FundingRate = fundingRateLong
		if fundingRateLong.IsZero() {
			pool.FundingRate = fundingRateShort.Neg()
		}

		totalLongOpenInterest := pool.GetTotalLongOpenInterest()
		totalShortOpenInterest := pool.GetTotalShortOpenInterest()

		blocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
		fundingAmountLong := types.CalcTakeAmount(totalLongOpenInterest, fundingRateLong).ToLegacyDec().QuoInt64(blocksPerYear)
		fundingAmountShort := types.CalcTakeAmount(totalShortOpenInterest, fundingRateShort).ToLegacyDec().QuoInt64(blocksPerYear)

		fundingShareLong := math.LegacyZeroDec()
		if totalShortOpenInterest.IsPositive() {
			fundingShareLong = fundingAmountLong.Quo(totalShortOpenInterest.ToLegacyDec())
		}

		fundingShareShort := math.LegacyZeroDec()
		if totalLongOpenInterest.IsPositive() {
			fundingShareShort = fundingAmountShort.Quo(totalLongOpenInterest.ToLegacyDec())
		}

		k.SetFundingRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.FundingRateBlock{
			BlockHeight:       ctx.BlockHeight(),
			BlockTime:         ctx.BlockTime().Unix(),
			FundingRateLong:   fundingRateLong,
			FundingRateShort:  fundingRateShort,
			FundingShareShort: fundingShareShort,
			FundingShareLong:  fundingShareLong,
		})
		k.SetPool(ctx, pool)
	}
}

func (k Keeper) ComputeFundingRate(ctx sdk.Context, pool types.Pool) (math.LegacyDec, math.LegacyDec) {
	// Custody amount for long is trading asset -
	// Liability amount for short is trading asset
	// popular_rate = fixed_rate * abs(Custody-Liability) / (Custody+Liability)
	totalLongOpenInterest := pool.GetTotalLongOpenInterest()
	totalShortOpenInterest := pool.GetTotalShortOpenInterest()

	if totalLongOpenInterest.IsZero() || totalShortOpenInterest.IsZero() {
		return math.LegacyZeroDec(), math.LegacyZeroDec()
	}

	if totalLongOpenInterest.Equal(totalShortOpenInterest) {
		return math.LegacyZeroDec(), math.LegacyZeroDec()
	}

	params := k.GetParams(ctx)
	if totalLongOpenInterest.GT(totalShortOpenInterest) {
		// long is popular
		// long pays short
		netLongRatio := (totalLongOpenInterest.Sub(totalShortOpenInterest)).ToLegacyDec().Quo((totalLongOpenInterest.Add(totalShortOpenInterest)).ToLegacyDec())
		fundingRate := netLongRatio.Mul(params.FixedFundingRate)
		if fundingRate.LT(params.MinimumFundingRate) {
			fundingRate = params.MinimumFundingRate
		}
		return fundingRate, math.LegacyZeroDec()
	} else {
		// short is popular
		// short pays long
		netShortRatio := (totalShortOpenInterest.Sub(totalLongOpenInterest)).ToLegacyDec().Quo((totalLongOpenInterest.Add(totalShortOpenInterest)).ToLegacyDec())
		fundingRate := netShortRatio.Mul(params.FixedFundingRate)
		if fundingRate.LT(params.MinimumFundingRate) {
			fundingRate = params.MinimumFundingRate
		}
		return math.LegacyZeroDec(), fundingRate
	}
}
