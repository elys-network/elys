package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (sdk.Int, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if len(params.StakeIncentives) < 1 || len(params.LpIncentives) < 1 {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNoNonInflationaryParams, "no inflationary params available")
	}

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroInt(), sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	baseCurrency := entry.Denom
	lpIncentive := params.LpIncentives[0]
	stkIncentive := params.StakeIncentives[0]

	if lpIncentive.NumEpochs < 1 || stkIncentive.NumEpochs < 1 {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNoNonInflationaryParams, "invalid inflationary params")
	}

	// Eden amount for LPs at a single epoch.
	edenAmountPerEpochLPs := lpIncentive.Amount.Quo(sdk.NewInt(lpIncentive.NumEpochs))

	// Eden amount for Stakers at a single epoch.
	edenAmountPerEpochStakers := stkIncentive.Amount.Quo(sdk.NewInt(stkIncentive.NumEpochs))

	if query.Denom == ptypes.Eden {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			epoch, found := k.epochsKeeper.GetEpochInfo(ctx, lpIncentive.EpochIdentifier)
			if !found {
				return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNoNonInflationaryParams, "no inflationary params available")
			}

			totalUSDCDeposit := k.bankKeeper.GetBalance(ctx, stabletypes.PoolAddress(), baseCurrency)
			if totalUSDCDeposit.Amount.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Calculate total Proxy TVL
			totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

			// Calculate stable stake pool share.
			poolShare := k.CalculatePoolShareForStableStakeLPs(ctx, totalProxyTVL, baseCurrency)

			// Epoch count in 24 hrs.
			epochDuration := int64(epoch.Duration.Seconds())

			// Eden amount for LP in 24hrs = Eden amount per a single epoch for LPs * number of epochs in a day
			edenAmountPerDay := edenAmountPerEpochLPs.Mul(sdk.NewInt(ptypes.SecondsPerDay)).Quo(sdk.NewInt(epochDuration))

			// Eden amount for stable stake LP in 24hrs
			edenAmountPerStableStakePerDay := sdk.NewDecFromInt(edenAmountPerDay).Mul(poolShare)

			// Calc Eden price in usdc
			edenPrice := k.EstimatePrice(ctx, sdk.NewCoin(ptypes.Eden, sdk.NewInt(1)), baseCurrency)

			// Eden Apr for usdc earn program = {(Eden allocated for stable stake pool per day*365*price{eden/usdc}/(total usdc deposit)}*100
			apr := edenAmountPerStableStakePerDay.MulInt(sdk.NewInt(100)).MulInt(sdk.NewInt(ptypes.DaysPerYear)).MulInt(edenPrice).QuoInt(totalUSDCDeposit.Amount)
			return apr.TruncateInt(), nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			epoch, found := k.epochsKeeper.GetEpochInfo(ctx, stkIncentive.EpochIdentifier)
			if !found {
				return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNoNonInflationaryParams, "no inflationary params available")
			}

			// Epoch count in 24 hrs.
			epochDuration := int64(epoch.Duration.Seconds())

			// Eden amount for LP in 24hrs = Eden amount per a single epoch for LPs * number of epochs in a day
			edenAmountPerDay := edenAmountPerEpochStakers.Mul(sdk.NewInt(ptypes.SecondsPerDay)).Quo(sdk.NewInt(epochDuration))

			// Update total committed states
			k.UpdateTotalCommitmentInfo(ctx, baseCurrency)
			totalStakedSnapshot := k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)

			// For Eden reward Apr for elys staking = {(amount of Eden allocated for staking per day)*365/( total elys staked + total Eden committed + total Eden boost committed)}*100
			apr := edenAmountPerDay.Mul(sdk.NewInt(100)).Mul(sdk.NewInt(ptypes.DaysPerYear)).Quo(totalStakedSnapshot)
			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			apr := params.InterestRate.MulInt(sdk.NewInt(100))
			return apr.TruncateInt(), nil
		} else {
			// Elys staking, Eden committed, EdenB committed.

			return sdk.ZeroInt(), nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := sdk.NewInt(lpIncentive.EdenBoostApr)
		return apr, nil
	}

	return sdk.ZeroInt(), nil
}
