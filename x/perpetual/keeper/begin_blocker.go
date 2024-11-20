package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
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
		fundingAmountLong := types.CalcTakeAmount(totalLongOpenInterest, fundingRateLong).ToLegacyDec().Quo(math.LegacyNewDec(blocksPerYear))
		fundingAmountShort := types.CalcTakeAmount(totalShortOpenInterest, fundingRateShort).ToLegacyDec().Quo(math.LegacyNewDec(blocksPerYear))

		k.SetFundingRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.FundingRateBlock{
			BlockHeight:        ctx.BlockHeight(),
			BlockTime:          ctx.BlockTime().Unix(),
			FundingRateLong:    fundingRateLong,
			FundingRateShort:   fundingRateShort,
			FundingAmountShort: fundingAmountShort,
			FundingAmountLong:  fundingAmountLong,
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

	fixedRate := k.GetParams(ctx).FixedFundingRate
	if totalLongOpenInterest.GT(totalShortOpenInterest) {
		// long is popular
		// long pays short
		netLongRatio := (totalLongOpenInterest.Sub(totalShortOpenInterest)).ToLegacyDec().Quo((totalLongOpenInterest.Add(totalShortOpenInterest)).ToLegacyDec())
		return netLongRatio.Mul(fixedRate), math.LegacyZeroDec()
	} else {
		// short is popular
		// short pays long
		netShortRatio := (totalShortOpenInterest.Sub(totalLongOpenInterest)).ToLegacyDec().Quo((totalLongOpenInterest.Add(totalShortOpenInterest)).ToLegacyDec())
		return math.LegacyZeroDec(), netShortRatio.Mul(fixedRate)
	}
}
