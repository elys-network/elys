package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	stabletypes "github.com/elys-network/elys/v5/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) CalculateStableStakeApr(ctx sdk.Context, query *types.QueryStableStakeAprRequest) (osmomath.BigDec, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return osmomath.ZeroBigDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// If we don't have enough params
	if query.Denom == ptypes.Eden {
		lpIncentive := params.LpIncentives
		if lpIncentive == nil || lpIncentive.EdenAmountPerYear.IsNil() {
			return osmomath.ZeroBigDec(), nil
		}

		totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
		if totalBlocksPerYear == 0 {
			return osmomath.ZeroBigDec(), nil
		}

		stableTvl := k.stableKeeper.TVL(ctx, stabletypes.UsdcPoolId)
		if stableTvl.IsZero() {
			return osmomath.ZeroBigDec(), nil
		}

		// Calculate total Proxy TVL
		totalProxyTVL, _ := k.CalculateProxyTVL(ctx, baseCurrency)

		edenAmount := lpIncentive.EdenAmountPerYear.Quo(sdkmath.NewInt(totalBlocksPerYear))

		edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, uint64(stabletypes.UsdcPoolId))
		if !found {
			return osmomath.ZeroBigDec(), nil
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := stableTvl.Mul(poolInfo.GetBigDecMultiplier())
		if totalProxyTVL.IsZero() {
			return osmomath.ZeroBigDec(), nil
		}
		stableStakePoolShare := proxyTVL.Quo(totalProxyTVL)

		stableStakeEdenAmount := osmomath.BigDecFromSDKInt(edenAmount).Mul(stableStakePoolShare)

		params := k.GetParams(ctx)
		poolMaxEdenAmount := params.GetBigDecMaxEdenRewardAprLps().
			Mul(proxyTVL).
			QuoInt64(totalBlocksPerYear).
			Quo(edenDenomPrice)
		stableStakeEdenAmount = osmomath.MinBigDec(stableStakeEdenAmount, poolMaxEdenAmount)

		// Eden Apr for usdc earn program
		apr := stableStakeEdenAmount.
			MulInt64(totalBlocksPerYear).
			Mul(edenDenomPrice).
			Quo(stableTvl)
		return apr, nil
	} else if query.Denom == ptypes.BaseCurrency {
		borrowPool, found := k.stableKeeper.GetPoolByDenom(ctx, baseCurrency)
		if !found {
			return osmomath.ZeroBigDec(), errorsmod.Wrap(types.ErrPoolNotFound, "pool not found")
		}
		res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
		if err != nil {
			return osmomath.ZeroBigDec(), err
		}
		apr := borrowPool.GetBigDecInterestRate().MulDec(res.BorrowRatio)
		return apr, nil
	}

	return osmomath.ZeroBigDec(), nil
}
