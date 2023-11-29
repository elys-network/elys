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

	if lpIncentive.TotalBlocksPerYear.IsZero() || stkIncentive.TotalBlocksPerYear.IsZero() {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNoNonInflationaryParams, "invalid inflationary params")
	}

	if query.Denom == ptypes.Eden {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			totalUSDCDeposit := k.bankKeeper.GetBalance(ctx, stabletypes.PoolAddress(), baseCurrency)
			if totalUSDCDeposit.Amount.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Calculate total Proxy TVL
			totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

			// Calculate stable stake pool share.
			poolShare := k.CalculatePoolShareForStableStakeLPs(ctx, totalProxyTVL, baseCurrency)

			// Eden amount for LP in 24hrs = AllocationEpochInBlocks is the number of block for 24 hrs
			edenAmountPerDay := lpIncentive.EdenAmountPerYear.Mul(lpIncentive.AllocationEpochInBlocks).Quo(lpIncentive.TotalBlocksPerYear)

			maxEdenAmountPerLps := params.MaxEdenRewardAprLps.Mul(totalProxyTVL).MulInt(lpIncentive.AllocationEpochInBlocks).QuoInt(lpIncentive.TotalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			edenAmountPerDay = sdk.MinInt(edenAmountPerDay, maxEdenAmountPerLps.TruncateInt())

			// Eden amount for stable stake LP in 24hrs
			edenAmountPerStableStakePerDay := sdk.NewDecFromInt(edenAmountPerDay).Mul(poolShare)

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenPrice := k.EstimatePrice(ctx, sdk.NewCoin(ptypes.Elys, sdk.NewInt(100000)), baseCurrency)

			// Eden Apr for usdc earn program = {(Eden allocated for stable stake pool per day*365*price{eden/usdc}/(total usdc deposit)}*100
			// we divide 100000 as we have use 100000elys as input in the price estimation
			apr := edenAmountPerStableStakePerDay.MulInt(sdk.NewInt(ptypes.DaysPerYear)).MulInt(edenPrice).MulInt(sdk.NewInt(100)).QuoInt(totalUSDCDeposit.Amount).QuoInt(sdk.NewInt(100000))
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
			edenAmountPerEpochStakersPerDay := stkIncentive.EdenAmountPerYear.Mul(stkIncentive.AllocationEpochInBlocks).Quo(stkIncentive.TotalBlocksPerYear)

			// Maximum eden based per distribution epoch on maximum APR - 30% by default
			// Allocated for staking per day = (0.3/365)* ( total elys staked + total Eden committed + total Eden boost committed)
			maxEdenAmountPerStakers := params.MaxEdenRewardAprStakers.MulInt(totalStakedSnapshot).MulInt(stkIncentive.AllocationEpochInBlocks).QuoInt(stkIncentive.TotalBlocksPerYear)

			// Use min amount (eden allocation from tokenomics and max apr based eden amount)
			edenAmountPerEpochStakersPerDay = sdk.MinInt(edenAmountPerEpochStakersPerDay, maxEdenAmountPerStakers.TruncateInt())

			// For Eden reward Apr for elys staking = {(amount of Eden allocated for staking per day)*365/( total elys staked + total Eden committed + total Eden boost committed)}*100
			apr := edenAmountPerEpochStakersPerDay.Mul(sdk.NewInt(ptypes.DaysPerYear)).Mul(sdk.NewInt(100)).Quo(totalStakedSnapshot)

			return apr, nil
		}
	} else if query.Denom == ptypes.BaseCurrency {
		if query.WithdrawType == commitmenttypes.EarnType_USDC_PROGRAM {
			params := k.stableKeeper.GetParams(ctx)
			apr := params.InterestRate.MulInt(sdk.NewInt(100))
			return apr.TruncateInt(), nil
		} else {
			// Elys staking, Eden committed, EdenB committed.
			params := k.GetParams(ctx)
			amt := params.DexRewardsStakers.Amount
			if amt.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// If no rewards were given.
			if params.DexRewardsStakers.NumBlocks.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// DexReward amount per day = amount distributed / duration(in seconds) * total seconds per day.
			// AllocationEpochInBlocks is the number of the block per day
			amtDexRewardPerDay := amt.MulInt(stkIncentive.AllocationEpochInBlocks).QuoInt(params.DexRewardsStakers.NumBlocks)

			// Calc Eden price in usdc
			// We put Elys as denom as Eden won't be avaialble in amm pool and has the same value as Elys
			edenPrice := k.EstimatePrice(ctx, sdk.NewCoin(ptypes.Elys, sdk.NewInt(10)), baseCurrency)

			// Update total committed states
			k.UpdateTotalCommitmentInfo(ctx, baseCurrency)
			totalStakedSnapshot := k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)

			// Ensure totalStakedSnapshot is not zero to avoid division by zero
			if totalStakedSnapshot.IsZero() {
				return sdk.ZeroInt(), nil
			}

			// Usdc apr for elys staking = (24 hour dex rewards in USDC generated for stakers) * 365*100/ {price ( elys/usdc)*( sum of (elys staked, Eden committed, Eden boost committed))}
			// we multiply 10 as we have use 10elys as input in the price estimation
			apr := amtDexRewardPerDay.MulInt(sdk.NewInt(ptypes.DaysPerYear)).MulInt(sdk.NewInt(100)).MulInt(sdk.NewInt(10)).QuoInt(edenPrice).QuoInt(totalStakedSnapshot)

			return apr.TruncateInt(), nil
		}
	} else if query.Denom == ptypes.EdenB {
		apr := lpIncentive.EdenBoostApr.MulInt(sdk.NewInt(100)).TruncateInt()
		return apr, nil
	}

	return sdk.ZeroInt(), nil
}
