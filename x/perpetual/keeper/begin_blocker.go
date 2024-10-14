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
				BlockTime:    ctx.BlockTime().Unix(),
			})

			err = k.UpdatePoolHealth(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}

			fundingRateLong, fundingRateShort := k.ComputeFundingRate(ctx, pool)

			k.SetFundingRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.FundingRateBlock{
				BlockHeight:      ctx.BlockHeight(),
				BlockTime:        ctx.BlockTime().Unix(),
				FundingRateLong:  fundingRateLong,
				FundingRateShort: fundingRateShort,
			})
		}
		k.SetPool(ctx, pool)
	}
}

func (k Keeper) ComputeFundingRate(ctx sdk.Context, pool types.Pool) (sdk.Dec, sdk.Dec) {
	// Custody amount for long is trading asset -
	// Liability amount for short is trading asset
	// popular_rate = fixed_rate * abs(Custody-Liability) / (Custody+Liability)
	totalCustodyLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		totalCustodyLong = totalCustodyLong.Add(asset.Custody)
	}

	totalLiabilitiesShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		totalLiabilitiesShort = totalLiabilitiesShort.Add(asset.Liabilities)
	}

	fixedRate := k.GetParams(ctx).FixedFundingRate
	if totalCustodyLong.GT(totalLiabilitiesShort) {
		// long is popular
		// long pays short
		if totalLiabilitiesShort.IsZero() {
			return sdk.ZeroDec(), sdk.ZeroDec()
		} else {
			netLongRatio := (totalCustodyLong.Sub(totalLiabilitiesShort)).ToLegacyDec().Quo((totalCustodyLong.Add(totalLiabilitiesShort)).ToLegacyDec())
			return netLongRatio.Mul(fixedRate), sdk.ZeroDec()
		}
	} else {
		// short is popular
		// short pays long
		if totalCustodyLong.IsZero() {
			return sdk.ZeroDec(), sdk.ZeroDec()
		} else {
			netShortRatio := (totalLiabilitiesShort.Sub(totalCustodyLong)).ToLegacyDec().Quo((totalCustodyLong.Add(totalLiabilitiesShort)).ToLegacyDec())
			return sdk.ZeroDec(), netShortRatio.Mul(fixedRate)
		}
	}
}
