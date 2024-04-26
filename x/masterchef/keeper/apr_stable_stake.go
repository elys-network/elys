package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateStableStakeApr(ctx sdk.Context, query *types.QueryStableStakeAprRequest) (math.Int, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if params.LpIncentives == nil {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroInt(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	baseCurrency := entry.Denom
	lpIncentive := params.LpIncentives

	if lpIncentive.TotalBlocksPerYear.IsZero() {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	if query.Denom == ptypes.Eden {
		stableTvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
		if stableTvl.IsZero() {
			return sdk.ZeroInt(), nil
		}

		// Calculate total Proxy TVL
		totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

		edenAmount := lpIncentive.EdenAmountPerYear.
			Quo(lpIncentive.TotalBlocksPerYear)

		edenDenomPrice := k.amm.GetEdenDenomPrice(ctx, baseCurrency)

		// Eden amount for stable stake LP in 24hrs
		stableStakePoolShare := k.CalculatePoolShareForStableStakeLPs(ctx, totalProxyTVL, baseCurrency)
		stableStakeEdenAmount := sdk.NewDecFromInt(edenAmount).Mul(stableStakePoolShare)

		params := k.GetParams(ctx)
		poolMaxEdenAmount := params.MaxEdenRewardAprLps.
			Mul(stableTvl).
			QuoInt(lpIncentive.TotalBlocksPerYear).
			Quo(edenDenomPrice)
		stableStakeEdenAmount = sdk.MinDec(stableStakeEdenAmount, poolMaxEdenAmount)

		// Eden Apr for usdc earn program = {(Eden allocated for stable stake pool per day*365*price{eden/usdc}/(total usdc deposit)}*100
		apr := stableStakeEdenAmount.
			MulInt(sdk.NewInt(ptypes.DaysPerYear)).
			Mul(edenDenomPrice).
			MulInt(sdk.NewInt(100)).
			Quo(stableTvl)
		return apr.TruncateInt(), nil
	} else if query.Denom == ptypes.BaseCurrency {
		params := k.stableKeeper.GetParams(ctx)
		res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
		if err != nil {
			return sdk.ZeroInt(), err
		}
		apr := params.InterestRate.Mul(res.BorrowRatio).MulInt(sdk.NewInt(100))
		return apr.TruncateInt(), nil
	}

	return sdk.ZeroInt(), nil
}
