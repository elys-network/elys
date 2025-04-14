package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateStableStakeApr(ctx sdk.Context, query *types.QueryStableStakeAprRequest) (sdkmath.LegacyDec, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if query.Denom == ptypes.Eden {
		lpIncentive := params.LpIncentives
		if lpIncentive == nil || lpIncentive.EdenAmountPerYear.IsNil() {
			return sdkmath.LegacyZeroDec(), nil
		}

		totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
		if totalBlocksPerYear == 0 {
			return sdkmath.LegacyZeroDec(), nil
		}

		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return sdkmath.LegacyZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
		}

		stableTvl := k.stableKeeper.TVL(ctx, stabletypes.UsdcPoolId)
		if stableTvl.IsZero() {
			return sdkmath.LegacyZeroDec(), nil
		}

		// Calculate total Proxy TVL
		totalProxyTVL, _ := k.CalculateProxyTVL(ctx, baseCurrency)

		edenAmount := lpIncentive.EdenAmountPerYear.Quo(sdkmath.NewInt(totalBlocksPerYear))

		edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, uint64(stabletypes.UsdcPoolId))
		if !found {
			return sdkmath.LegacyZeroDec(), nil
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := stableTvl.Mul(poolInfo.Multiplier)
		if totalProxyTVL.IsZero() {
			return sdkmath.LegacyZeroDec(), nil
		}
		stableStakePoolShare := proxyTVL.Quo(totalProxyTVL)

		stableStakeEdenAmount := sdkmath.LegacyNewDecFromInt(edenAmount).Mul(stableStakePoolShare)

		params := k.GetParams(ctx)
		poolMaxEdenAmount := params.MaxEdenRewardAprLps.
			Mul(proxyTVL).
			QuoInt64(totalBlocksPerYear).
			Quo(edenDenomPrice)
		stableStakeEdenAmount = sdkmath.LegacyMinDec(stableStakeEdenAmount, poolMaxEdenAmount)

		// Eden Apr for usdc earn program
		apr := stableStakeEdenAmount.
			MulInt64(totalBlocksPerYear).
			Mul(edenDenomPrice).
			Quo(stableTvl)
		return apr, nil
	} else if query.Denom == ptypes.BaseCurrency {
		borrowPool, found := k.stableKeeper.GetPoolByDenom(ctx, query.Denom)
		if !found {
			return math.LegacyZeroDec(), errorsmod.Wrap(types.ErrPoolNotFound, "pool not found")
		}
		res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
		if err != nil {
			return sdkmath.LegacyZeroDec(), err
		}
		apr := borrowPool.InterestRate.Mul(res.BorrowRatio)
		return apr, nil
	}

	return sdkmath.LegacyZeroDec(), nil
}
