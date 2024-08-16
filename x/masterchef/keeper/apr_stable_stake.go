package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateStableStakeApr(ctx sdk.Context, query *types.QueryStableStakeAprRequest) (sdk.Dec, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if query.Denom == ptypes.Eden {
		lpIncentive := params.LpIncentives
		if lpIncentive == nil || lpIncentive.EdenAmountPerYear.IsNil() {
			return sdk.ZeroDec(), nil
		}

		totalBlocksPerYear := k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear
		if totalBlocksPerYear == 0 {
			return sdk.ZeroDec(), nil
		}

		baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if !found {
			return sdk.ZeroDec(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
		}

		stableTvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
		if stableTvl.IsZero() {
			return sdk.ZeroDec(), nil
		}

		// Calculate total Proxy TVL
		totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

		edenAmount := lpIncentive.EdenAmountPerYear.Quo(sdk.NewInt(totalBlocksPerYear))

		edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, uint64(stabletypes.PoolId))
		if !found {
			return sdk.ZeroDec(), nil
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := stableTvl.Mul(poolInfo.Multiplier)
		if totalProxyTVL.IsZero() {
			return sdk.ZeroDec(), nil
		}
		stableStakePoolShare := proxyTVL.Quo(totalProxyTVL)

		stableStakeEdenAmount := sdk.NewDecFromInt(edenAmount).Mul(stableStakePoolShare)

		params := k.GetParams(ctx)
		poolMaxEdenAmount := params.MaxEdenRewardAprLps.
			Mul(proxyTVL).
			QuoInt64(totalBlocksPerYear).
			Quo(edenDenomPrice)
		stableStakeEdenAmount = sdk.MinDec(stableStakeEdenAmount, poolMaxEdenAmount)

		// Eden Apr for usdc earn program
		apr := stableStakeEdenAmount.
			MulInt64(totalBlocksPerYear).
			Mul(edenDenomPrice).
			Quo(stableTvl)
		return apr, nil
	} else if query.Denom == ptypes.BaseCurrency {
		params := k.stableKeeper.GetParams(ctx)
		res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
		if err != nil {
			return sdk.ZeroDec(), err
		}
		apr := params.InterestRate.Mul(res.BorrowRatio)
		return apr, nil
	}

	return sdk.ZeroDec(), nil
}
