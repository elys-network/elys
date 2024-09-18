package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	currentHeight := ctx.BlockHeight()
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
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
			})

			err = k.UpdatePoolHealth(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
			err = k.UpdateFundingRate(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}

			// account custody from long position
			totalCustodyLong := sdk.ZeroInt()
			for _, asset := range pool.PoolAssetsLong {
				totalCustodyLong = totalCustodyLong.Add(asset.Custody)
			}

			// account custody from short position
			totalCustodyShort := sdk.ZeroInt()
			for _, asset := range pool.PoolAssetsShort {
				totalCustodyShort = totalCustodyShort.Add(asset.Custody)
			}

			fundingAmountLong := types.CalcTakeAmount(totalCustodyLong, pool.FundingRate)
			fundingAmountShort := sdk.ZeroInt()

			fundingRateLong := pool.FundingRate
			fundingRateShort := sdk.ZeroDec()

			// if funding rate is negative, collect from short position
			if pool.FundingRate.IsNegative() {
				fundingAmountShort = types.CalcTakeAmount(totalCustodyShort, pool.FundingRate)
				fundingAmountLong = sdk.ZeroInt()

				fundingRateLong = sdk.ZeroDec()
				fundingRateShort = pool.FundingRate
			}
			k.SetFundingRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.FundingRateBlock{
				FundingRate:        pool.FundingRate,
				BlockHeight:        ctx.BlockHeight(),
				FundingAmountShort: fundingAmountShort,
				FundingAmountLong:  fundingAmountLong,
				FundingRateLong:    fundingRateLong,
				FundingRateShort:   fundingRateShort,
			})

			err = k.UpdateNetOpenInterest(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
		}

		k.SetPool(ctx, pool)
	}
}
