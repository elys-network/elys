package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) CalculateApr(ctx sdk.Context, query *types.QueryAprRequest) (math.Int, error) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update params
	defer k.SetParams(ctx, params)

	// If we don't have enough params
	if params.StakeIncentives == nil || params.LpIncentives == nil {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrNoInflationaryParams, "no inflationary params available")
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroInt(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	baseCurrency := entry.Denom
	lpIncentive := params.LpIncentives
	stkIncentive := params.StakeIncentives

	if lpIncentive.TotalBlocksPerYear.IsZero() || stkIncentive.TotalBlocksPerYear.IsZero() {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	if query.Denom == ptypes.Eden {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			stableTvl := k.stableKeeper.TVL(ctx, k.oracleKeeper, baseCurrency)
			if stableTvl.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Calculate total Proxy TVL
			totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

			// Eden amount for LP in 24hrs = EpochNumBlocks is the number of block for 24 hrs
			epochEdenAmount := lpIncentive.EdenAmountPerYear.
				Mul(lpIncentive.EpochNumBlocks).
				Quo(lpIncentive.TotalBlocksPerYear)

			edenDenomPrice := k.GetEdenDenomPrice(ctx, baseCurrency)
			epochLpsMaxEdenAmount := params.MaxEdenRewardAprLps.
				Mul(totalProxyTVL).
				MulInt(lpIncentive.EpochNumBlocks).
				QuoInt(lpIncentive.TotalBlocksPerYear).
				Quo(edenDenomPrice)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			epochEdenAmount = sdk.MinInt(epochEdenAmount, epochLpsMaxEdenAmount.TruncateInt())

			// Eden amount for stable stake LP in 24hrs
			stableStakePoolShare := k.CalculatePoolShareForStableStakeLPs(ctx, totalProxyTVL, baseCurrency)
			epochStableStakeEdenAmount := sdk.NewDecFromInt(epochEdenAmount).Mul(stableStakePoolShare)

			// Eden Apr for usdc earn program = {(Eden allocated for stable stake pool per day*365*price{eden/usdc}/(total usdc deposit)}*100
			apr := epochStableStakeEdenAmount.
				MulInt(sdk.NewInt(ptypes.DaysPerYear)).
				Mul(edenDenomPrice).
				MulInt(sdk.NewInt(100)).
				Quo(stableTvl)
			return apr.TruncateInt(), nil
		} else {
			// Elys staking, Eden committed, EdenB committed.

			// Update total committed states
			k.UpdateTotalCommitmentInfo(ctx, baseCurrency)
			totalStakedSnapshot := k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Calculate
			epochStakersEdenAmount := stkIncentive.EdenAmountPerYear.
				Mul(stkIncentive.EpochNumBlocks).
				Quo(stkIncentive.TotalBlocksPerYear)

			// Maximum eden based per distribution epoch on maximum APR - 30% by default
			// Allocated for staking per day = (0.3/365)* ( total elys staked + total Eden committed + total Eden boost committed)
			epochStakersMaxEdenAmount := params.MaxEdenRewardAprStakers.
				MulInt(totalStakedSnapshot).
				MulInt(stkIncentive.EpochNumBlocks).
				QuoInt(stkIncentive.TotalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			epochStakersEdenAmount = sdk.MinInt(epochStakersEdenAmount, epochStakersMaxEdenAmount.TruncateInt())

			// For Eden reward Apr for elys staking = {(amount of Eden allocated for staking per day)*365/( total elys staked + total Eden committed + total Eden boost committed)}*100
			apr := epochStakersEdenAmount.
				Mul(sdk.NewInt(ptypes.DaysPerYear)).
				Mul(sdk.NewInt(100)).
				Quo(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			res, err := k.stableKeeper.BorrowRatio(ctx, &stabletypes.QueryBorrowRatioRequest{})
			if err != nil {
				return sdk.ZeroInt(), err
			}
			apr := params.InterestRate.Mul(res.BorrowRatio).MulInt(sdk.NewInt(100))
			return apr.TruncateInt(), nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			params := k.GetParams(ctx)
			amount := params.DexRewardsStakers.Amount
			if amount.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// If no rewards were given.
			if params.DexRewardsStakers.NumBlocks.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenPrice := k.EstimatePrice(ctx, ptypes.Elys, baseCurrency)
			if edenPrice.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Update total committed states
			k.UpdateTotalCommitmentInfo(ctx, baseCurrency)
			totalStakedSnapshot := k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// DexReward amount per day = amount distributed / duration(in seconds) * total seconds per day.
			// EpochNumBlocks is the number of the block per day
			dailyDexRewardAmount := amount.MulInt(stkIncentive.EpochNumBlocks).QuoInt(params.DexRewardsStakers.NumBlocks)

			// Usdc apr for elys staking = (24 hour dex rewards in USDC generated for stakers) * 365*100/ {price ( elys/usdc)*( sum of (elys staked, Eden committed, Eden boost committed))}
			// we multiply 10 as we have use 10elys as input in the price estimation
			apr := dailyDexRewardAmount.
				MulInt(sdk.NewInt(ptypes.DaysPerYear)).
				MulInt(sdk.NewInt(100)).
				Quo(edenPrice).
				QuoInt(totalStakedSnapshot)

			return apr.TruncateInt(), nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := types.EdenBoostApr.MulInt(sdk.NewInt(100)).TruncateInt()
		return apr, nil
	}

	return sdk.ZeroInt(), nil
}
