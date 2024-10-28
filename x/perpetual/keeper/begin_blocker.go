package keeper

import (
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

		blocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear
		fundingAmountLong := types.CalcTakeAmount(totalCustodyLong, fundingRateLong).ToLegacyDec().Quo(sdk.NewDec(blocksPerYear))
		fundingAmountShort := types.CalcTakeAmount(totalCustodyShort, fundingRateShort).ToLegacyDec().Quo(sdk.NewDec(blocksPerYear))

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

func (k Keeper) ComputeFundingRate(ctx sdk.Context, pool types.Pool) (sdk.Dec, sdk.Dec) {
	// Custody amount for long is trading asset -
	// Liability amount for short is trading asset
	// popular_rate = fixed_rate * abs(Custody-Liability) / (Custody+Liability)
	totalCustodyLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		// We subtract asset.Collateral from totalCustodyLong because for long with collateral same as trading asset and user will
		// be charged for that the collateral as well even though they have already given that amount to the pool.
		// For LONG, asset.Custody will be 0 only for base currency but asset.Collateral won't be zero for base currency and trading asset
		// We subtract asset.Collateral only when asset is trading asset and in that case asset.Custody won't be zero
		// For base currency, asset.Collateral might not be 0 but asset.Custody will be 0 in LONG
		// !asset.Custody.IsZero() ensures that asset is trading asset for LONG
		if !asset.Custody.IsZero() {
			totalCustodyLong = totalCustodyLong.Add(asset.Custody).Sub(asset.Collateral)
		}
	}

	totalLiabilitiesShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		totalLiabilitiesShort = totalLiabilitiesShort.Add(asset.Liabilities)
	}

	if totalCustodyLong.IsZero() || totalLiabilitiesShort.IsZero() {
		return sdk.ZeroDec(), sdk.ZeroDec()
	}

	fixedRate := k.GetParams(ctx).FixedFundingRate
	if totalCustodyLong.GT(totalLiabilitiesShort) {
		// long is popular
		// long pays short
		netLongRatio := (totalCustodyLong.Sub(totalLiabilitiesShort)).ToLegacyDec().Quo((totalCustodyLong.Add(totalLiabilitiesShort)).ToLegacyDec())
		return netLongRatio.Mul(fixedRate), sdk.ZeroDec()
	} else {
		// short is popular
		// short pays long
		netShortRatio := (totalLiabilitiesShort.Sub(totalCustodyLong)).ToLegacyDec().Quo((totalCustodyLong.Add(totalLiabilitiesShort)).ToLegacyDec())
		return sdk.ZeroDec(), netShortRatio.Mul(fixedRate)
	}
}
