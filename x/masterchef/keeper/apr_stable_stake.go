package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateStableStakeApr(ctx sdk.Context, query *types.QueryStableStakeAprRequest) (elystypes.Dec34, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if query.Denom == ptypes.Eden {
		lpIncentive := params.LpIncentives
		if lpIncentive == nil || lpIncentive.EdenAmountPerYear.IsNil() {
			return elystypes.ZeroDec34(), nil
		}

		totalBlocksPerYear := int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
		if totalBlocksPerYear == 0 {
			return elystypes.ZeroDec34(), nil
		}

		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return elystypes.ZeroDec34(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
		}

		stableTvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
		if stableTvl.IsZero() {
			return elystypes.ZeroDec34(), nil
		}

		// Calculate total Proxy TVL
		totalProxyTVL, _ := k.CalculateProxyTVL(ctx, baseCurrency)

		edenAmount := lpIncentive.EdenAmountPerYear.Quo(sdkmath.NewInt(totalBlocksPerYear))

		edenDenomPrice, decimals := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, uint64(stabletypes.PoolId))
		if !found {
			return elystypes.ZeroDec34(), nil
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := stableTvl.MulLegacyDec(poolInfo.Multiplier)
		if totalProxyTVL.IsZero() {
			return elystypes.ZeroDec34(), nil
		}
		stableStakePoolShare := proxyTVL.Quo(totalProxyTVL)

		stableStakeEdenAmount := stableStakePoolShare.MulInt(edenAmount)

		params := k.GetParams(ctx)
		poolMaxEdenAmount := proxyTVL.
			MulLegacyDec(params.MaxEdenRewardAprLps).
			QuoInt64(totalBlocksPerYear).
			Quo(edenDenomPrice.QuoInt(ammtypes.OneTokenUnit(decimals)))
		stableStakeEdenAmount = elystypes.MinDec34(stableStakeEdenAmount, poolMaxEdenAmount)

		// Eden Apr for usdc earn program
		apr := stableStakeEdenAmount.
			MulInt64(totalBlocksPerYear).
			Mul(edenDenomPrice).
			QuoInt(ammtypes.OneTokenUnit(decimals)).
			Quo(stableTvl)
		return apr, nil
	} else if query.Denom == ptypes.BaseCurrency {
		params := k.stableKeeper.GetParams(ctx)
		res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
		if err != nil {
			return elystypes.ZeroDec34(), err
		}
		apr := elystypes.NewDec34FromLegacyDec(params.InterestRate).MulLegacyDec(res.BorrowRatio)
		return apr, nil
	}

	return elystypes.ZeroDec34(), nil
}
