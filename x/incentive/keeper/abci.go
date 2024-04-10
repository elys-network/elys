package keeper

import (
	"errors"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// EndBlocker of incentive module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	// Burn EdenB tokens if staking changed
	k.BurnEdenBIfElysStakingReduced(ctx)

	// // Rewards distribution
	k.ProcessRewardsDistribution(ctx)
}

func (k Keeper) TakeDelegationSnapshot(ctx sdk.Context, addr string) {
	// Calculate delegated amount per delegator
	delAmount := k.CalcDelegationAmount(ctx, addr)

	elysStaked := types.ElysStaked{
		Address: addr,
		Amount:  delAmount,
	}

	// Set Elys staked amount
	k.SetElysStaked(ctx, elysStaked)
}

func (k Keeper) BurnEdenBIfElysStakingReduced(ctx sdk.Context) {
	addrs := k.GetAllElysStakeChange(ctx)

	// Handle addresses recorded on AfterDelegationModified
	// This hook is exposed for genesis delegations as well
	for _, delAddr := range addrs {
		k.BurnEdenBFromElysUnstaking(ctx, delAddr)
		k.TakeDelegationSnapshot(ctx, delAddr.String())
		k.RemoveElysStakeChange(ctx, delAddr)
	}
}

// Rewards distribution
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) {
	// Read tokenomics time based inflation params and update incentive module params.
	if !k.ProcessUpdateIncentiveParams(ctx) {
		ctx.Logger().Error("Invalid tokenomics params", "error", errors.New("invalid tokenomics params"))
		return
	}

	stakerEpoch, stakeIncentive := k.IsStakerRewardsDistributionEpoch(ctx)
	if stakerEpoch {
		err := k.UpdateStakersRewardsUnclaimed(ctx, *stakeIncentive)
		if err != nil {
			ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
		}
	}

	lpsEpoch, lpIncentive := k.IsLPRewardsDistributionEpoch(ctx)
	if lpsEpoch {
		err := k.UpdateLPRewardsUnclaimed(ctx, *lpIncentive)
		if err != nil {
			ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
		}
	}
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) bool {
	// Non-linear inflation per year happens and this includes yearly inflation data
	listTimeBasedInflations := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflations) < 1 {
		return false
	}

	params := k.GetParams(ctx)

	for _, inflation := range listTimeBasedInflations {
		// Finding only current inflation data - and skip rest
		if inflation.StartBlockHeight > uint64(ctx.BlockHeight()) || inflation.EndBlockHeight < uint64(ctx.BlockHeight()) {
			continue
		}

		totalBlocksPerYear := sdk.NewInt(int64(inflation.EndBlockHeight - inflation.StartBlockHeight + 1))

		// ------------- LP Incentive parameter -------------
		totalDistributionEpochPerYear := totalBlocksPerYear
		// If totalDistributionEpochPerYear is zero, we skip this inflation to avoid division by zero
		if totalBlocksPerYear == sdk.ZeroInt() {
			continue
		}
		currentEpochInBlocks := sdk.NewInt(ctx.BlockHeight() - int64(inflation.StartBlockHeight)).
			Mul(totalDistributionEpochPerYear).
			Quo(totalBlocksPerYear)

		incentiveInfo := types.IncentiveInfo{
			// reward amount in eden for 1 year
			EdenAmountPerYear: sdk.NewInt(int64(inflation.Inflation.LmRewards)),
			// starting block height of the distribution
			DistributionStartBlock: sdk.NewInt(int64(inflation.StartBlockHeight)),
			// distribution duration - block number per year
			TotalBlocksPerYear: totalBlocksPerYear,
			// current epoch in block number
			CurrentEpochInBlocks: currentEpochInBlocks,
		}

		if params.LpIncentives == nil {
			params.LpIncentives = &incentiveInfo
		} else {
			// If any of block number related parameter changed, we re-calculate the current epoch
			if params.LpIncentives.DistributionStartBlock != incentiveInfo.DistributionStartBlock ||
				params.LpIncentives.TotalBlocksPerYear != incentiveInfo.TotalBlocksPerYear {
				params.LpIncentives.CurrentEpochInBlocks = currentEpochInBlocks
			}
			params.LpIncentives.EdenAmountPerYear = incentiveInfo.EdenAmountPerYear
			params.LpIncentives.DistributionStartBlock = incentiveInfo.DistributionStartBlock
			params.LpIncentives.TotalBlocksPerYear = incentiveInfo.TotalBlocksPerYear
		}

		// ------------- Stakers parameter -------------
		totalDistributionEpochPerYear = totalBlocksPerYear
		currentEpochInBlocks = sdk.NewInt(ctx.BlockHeight() - int64(inflation.StartBlockHeight)).Mul(totalDistributionEpochPerYear).Quo(totalBlocksPerYear)
		incentiveInfo = types.IncentiveInfo{
			// reward amount in eden for 1 year
			EdenAmountPerYear: sdk.NewInt(int64(inflation.Inflation.IcsStakingRewards)),
			// starting block height of the distribution
			DistributionStartBlock: sdk.NewInt(int64(inflation.StartBlockHeight)),
			// distribution duration - block number per year
			TotalBlocksPerYear: totalBlocksPerYear,
			// current epoch in block number
			CurrentEpochInBlocks: currentEpochInBlocks,
		}

		if params.StakeIncentives == nil {
			params.StakeIncentives = &incentiveInfo
		} else {
			// If any of block number related parameter changed, we re-calculate the current epoch
			if params.StakeIncentives.DistributionStartBlock != incentiveInfo.DistributionStartBlock ||
				params.StakeIncentives.TotalBlocksPerYear != incentiveInfo.TotalBlocksPerYear {
				params.StakeIncentives.CurrentEpochInBlocks = currentEpochInBlocks
			}
			params.StakeIncentives.EdenAmountPerYear = incentiveInfo.EdenAmountPerYear
			params.StakeIncentives.DistributionStartBlock = incentiveInfo.DistributionStartBlock
			params.StakeIncentives.TotalBlocksPerYear = incentiveInfo.TotalBlocksPerYear
		}
		break
	}

	k.SetParams(ctx, params)
	return true
}

func (k Keeper) IsStakerRewardsDistributionEpoch(ctx sdk.Context) (bool, *types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false, nil
	}

	// If we don't have enough params
	if params.StakeIncentives == nil {
		return false, nil
	}

	// Incentive params initialize
	stakeIncentive := params.StakeIncentives

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())
	if stakeIncentive.DistributionStartBlock.GT(curBlockHeight) {
		return false, nil
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.CurrentEpochInBlocks = stakeIncentive.CurrentEpochInBlocks
	if stakeIncentive.CurrentEpochInBlocks.GTE(stakeIncentive.TotalBlocksPerYear) || curBlockHeight.GT(stakeIncentive.TotalBlocksPerYear.Add(stakeIncentive.DistributionStartBlock)) {
		params.StakeIncentives = nil
		k.SetParams(ctx, params)
		return false, nil
	}

	params.StakeIncentives.CurrentEpochInBlocks = stakeIncentive.CurrentEpochInBlocks
	k.SetParams(ctx, params)

	// return found, stake incentive params
	return true, stakeIncentive
}

func (k Keeper) IsLPRewardsDistributionEpoch(ctx sdk.Context) (bool, *types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false, nil
	}

	// If we don't have enough params
	if params.LpIncentives == nil {
		return false, nil
	}

	// Incentive params initialize
	lpIncentive := params.LpIncentives

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())
	if lpIncentive.DistributionStartBlock.GT(curBlockHeight) {
		return false, nil
	}

	// Increase current epoch of Stake incentive param
	lpIncentive.CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks.Add(sdk.NewInt(1))
	if lpIncentive.CurrentEpochInBlocks.GTE(lpIncentive.TotalBlocksPerYear) || curBlockHeight.GT(lpIncentive.TotalBlocksPerYear.Add(lpIncentive.DistributionStartBlock)) {
		params.LpIncentives = nil
		k.SetParams(ctx, params)
		return false, nil
	}

	params.LpIncentives.CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks
	k.SetParams(ctx, params)

	// return found, lp incentive params
	return true, lpIncentive
}

// Update unclaimed token amount
// Called back through epoch hook
func (k Keeper) UpdateStakersRewardsUnclaimed(ctx sdk.Context, stakeIncentive types.IncentiveInfo) error {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx, baseCurrency)

	// Collect DEX revenue while tracking 65% of it for LPs reward calculation
	// Assume these are collected in USDC
	_, dexRevenueForLpsPerDistribution, dexRevenueForStakersPerDistribution := k.CollectDEXRevenue(ctx)

	// Calculate each portion of Gas fees collected - stakers, LPs
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(k.tci.TotalFeesCollected...)
	rewardPortionForLps := k.GetDEXRewardPortionForLPs(ctx)
	rewardPortionForStakers := k.GetDEXRewardPortionForStakers(ctx)
	gasFeesForLps := gasFeeCollectedDec.MulDecTruncate(rewardPortionForLps)
	gasFeesForStakers := gasFeeCollectedDec.MulDecTruncate(rewardPortionForStakers)

	// Sum Dex revenue for stakers + Gas fees for stakers and name it dex Revenus for stakers
	// But won't sum dex revenue for LPs and gas fees for LPs as the LP revenue will be rewared by pool.
	dexRevenueForStakersPerDistribution = dexRevenueForStakersPerDistribution.Add(gasFeesForStakers...)

	// USDC amount in sdk.Dec type
	dexRevenueLPsAmtPerDistribution := dexRevenueForLpsPerDistribution.AmountOf(baseCurrency)
	dexRevenueStakersAmtPerDistribution := dexRevenueForStakersPerDistribution.AmountOf(baseCurrency)
	gasFeesLPsAmtPerDistribution := gasFeesForLps.AmountOf(baseCurrency)

	// Calculate eden amount per epoch
	params := k.GetParams(ctx)

	// Ensure stakeIncentive.TotalBlocksPerYear or stakeIncentive.EpochNumBlocks are not zero to avoid division by zero
	if stakeIncentive.TotalBlocksPerYear.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate
	epochStakersEdenAmount := stakeIncentive.EdenAmountPerYear.
		Quo(stakeIncentive.TotalBlocksPerYear)

	// Maximum eden based per distribution epoch on maximum APR - 30% by default
	// Allocated for staking per day = (0.3/365)* ( total elys staked + total Eden committed + total Eden boost committed)
	epochStakersMaxEdenAmount := params.MaxEdenRewardAprStakers.
		MulInt(k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)).
		QuoInt(stakeIncentive.TotalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	epochStakersEdenAmount = sdk.MinInt(epochStakersEdenAmount, epochStakersMaxEdenAmount.TruncateInt())

	// Track the DEX rewards distribution for stakers
	// Add dexRevenue amount that was tracked by Lp tracker
	dexRevenueStakersAmtPerDistribution = dexRevenueStakersAmtPerDistribution.Add(params.DexRewardsStakers.AmountCollectedByOtherTracker)
	// Incrase total dex rewards given
	params.DexRewardsStakers.Amount = dexRevenueStakersAmtPerDistribution
	// Reset amount from other tracker
	params.DexRewardsStakers.AmountCollectedByOtherTracker = sdk.ZeroDec()
	// Don't increase Lps rewards blocks, it will be increased whenever LP distribution epoch happens.
	params.DexRewardsLps.AmountCollectedByOtherTracker = dexRevenueLPsAmtPerDistribution.
		Add(gasFeesLPsAmtPerDistribution)
	k.SetParams(ctx, params)

	totalEdenGiven := sdk.ZeroInt()
	totalRewardsGiven := sdk.ZeroInt()
	// Process to increase uncomitted token amount of Eden & Eden boost
	k.cmk.IterateCommitments(
		ctx, func(commitments ctypes.Commitments) bool {
			// Commitment owner
			creator := commitments.Creator
			_, err := sdk.AccAddressFromBech32(creator)
			if err != nil {
				// This could be validator address
				return false
			}

			rewardsByElysStaking := sdk.NewCoins()
			rewardsByEdenCommitted := sdk.NewCoins()
			rewardsByEdenBCommitted := sdk.NewCoins()

			newSumEdenRewardsUnClaimed := sdk.ZeroInt()
			newSumEdenBRewardsUnClaimed := sdk.ZeroInt()
			newSumDexRewardsUnClaimed := sdk.ZeroInt()

			// Calculate delegated amount per delegator
			bondedDelAmount := k.CalcBondedDelegationAmount(ctx, creator)

			// Calculate new unclaimed Eden tokens from Eden & Eden boost committed, Dex rewards distribution
			// Distribute gas fees to stakers

			// Calculate new unclaimed Eden tokens from Elys staked
			// ----------------------------------------------------------
			newUnclaimedEdenTokens, dexRewards, dexRewardsByStakers := k.CalcRewardsForStakersByElysStaked(
				ctx, bondedDelAmount, epochStakersEdenAmount, dexRevenueStakersAmtPerDistribution,
			)
			_ = dexRewardsByStakers

			// Total
			totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Sum for each loop
			newSumEdenRewardsUnClaimed = newSumEdenRewardsUnClaimed.Add(newUnclaimedEdenTokens)
			newSumDexRewardsUnClaimed = newSumDexRewardsUnClaimed.Add(dexRewards)

			// Store Eden rewards by Elys staking
			newEdenFromElysStaking := newUnclaimedEdenTokens
			newDexRewardFromElysStaking := dexRewards
			// ----------------------------------------------------------

			// Calculate new unclaimed Eden tokens from Eden committed
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)
			newUnclaimedEdenTokens, dexRewards = k.CalcRewardsForStakersByCommitted(
				ctx, edenCommitted, epochStakersEdenAmount, dexRevenueStakersAmtPerDistribution,
			)

			// Total
			totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Sum for each loop
			newSumEdenRewardsUnClaimed = newSumEdenRewardsUnClaimed.Add(newUnclaimedEdenTokens)
			newSumDexRewardsUnClaimed = newSumDexRewardsUnClaimed.Add(dexRewards)

			// Sub bucket
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(baseCurrency, dexRewards))
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// Calculate new unclaimed Eden tokens from Eden Boost committed
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			edenBoostCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
			newUnclaimedEdenTokens, dexRewards = k.CalcRewardsForStakersByCommitted(
				ctx,
				edenBoostCommitted,
				epochStakersEdenAmount,
				dexRevenueStakersAmtPerDistribution,
			)

			// Total
			totalEdenGiven = totalEdenGiven.Add(newUnclaimedEdenTokens)
			totalRewardsGiven = totalRewardsGiven.Add(dexRewards)

			// Sum for each loop
			newSumEdenRewardsUnClaimed = newSumEdenRewardsUnClaimed.Add(newUnclaimedEdenTokens)
			newSumDexRewardsUnClaimed = newSumDexRewardsUnClaimed.Add(dexRewards)

			// Sub bucket
			rewardsByEdenBCommitted = rewardsByEdenBCommitted.Add(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
			rewardsByEdenBCommitted = rewardsByEdenBCommitted.Add(sdk.NewCoin(baseCurrency, dexRewards))
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// Add Eden rewards from Elys staking
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.Eden, newEdenFromElysStaking))
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(baseCurrency, newDexRewardFromElysStaking))
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// Calculate new unclaimed Eden-Boost tokens for staker and Eden token holders
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			newEdenBTokens, newEdenBFromElysStaking, newEdenBFromEdenCommited := k.CalculateEdenBoostRewards(
				ctx, bondedDelAmount, commitments, stakeIncentive, types.EdenBoostApr)
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.EdenB, newEdenBFromElysStaking))
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(ptypes.EdenB, newEdenBFromEdenCommited))

			newSumEdenBRewardsUnClaimed = newSumEdenBRewardsUnClaimed.Add(newEdenBTokens)
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			return false
		},
	)

	// Calcualte the remainings
	edenRemained := epochStakersEdenAmount.Sub(totalEdenGiven)
	dexRewardsRemained := dexRevenueStakersAmtPerDistribution.Sub(sdk.NewDecFromInt(totalRewardsGiven))

	// if edenRemained is negative, override it with zero
	if edenRemained.IsNegative() {
		edenRemained = sdk.ZeroInt()
	}
	// if dexRewardsRemained is negative, override it with zero
	if dexRewardsRemained.IsNegative() {
		dexRewardsRemained = sdk.ZeroDec()
	}

	// Fund community the remain coins
	// ----------------------------------
	edenRemainedCoin := sdk.NewDecCoin(ptypes.Eden, edenRemained)
	dexRewardsRemainedCoin := sdk.NewDecCoinFromDec(baseCurrency, dexRewardsRemained)

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(edenRemainedCoin)
	feePool.CommunityPool = feePool.CommunityPool.Add(dexRewardsRemainedCoin)
	k.SetFeePool(ctx, feePool)
	// ----------------------------------

	return nil
}

// Update unclaimed token amount
// Called back through epoch hook
func (k Keeper) UpdateLPRewardsUnclaimed(ctx sdk.Context, lpIncentive types.IncentiveInfo) error {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	params := k.GetParams(ctx)

	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx, baseCurrency)

	// Collect DEX revenue while tracking 65% of it for LPs reward calculation
	// Assume these are collected in USDC
	_, dexRevenueForLpsPerDistribution, dexRevenueForStakersPerDistribution := k.CollectDEXRevenue(ctx)

	// Calculate each portion of Gas fees collected - stakers, LPs
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(k.tci.TotalFeesCollected...)
	rewardPortionForLps := k.GetDEXRewardPortionForLPs(ctx)
	rewardPortionForStakers := k.GetDEXRewardPortionForStakers(ctx)
	gasFeesForLpsPerDistribution := gasFeeCollectedDec.MulDecTruncate(rewardPortionForLps)
	gasFeesForStakersPerDistribution := gasFeeCollectedDec.MulDecTruncate(rewardPortionForStakers)

	// Sum Dex revenue for stakers + Gas fees for stakers and name it dex Revenus for stakers
	// But won't sum dex revenue for LPs and gas fees for LPs as the LP revenue will be rewared by pool.
	dexRevenueForStakersPerDistribution = dexRevenueForStakersPerDistribution.Add(gasFeesForStakersPerDistribution...)

	// USDC amount in sdk.Dec type
	dexRevenueLPsAmtPerDistribution := dexRevenueForLpsPerDistribution.AmountOf(baseCurrency)
	dexRevenueStakersAmtPerDistribution := dexRevenueForStakersPerDistribution.AmountOf(baseCurrency)
	gasFeesLPsAmtPerDistribution := gasFeesForLpsPerDistribution.AmountOf(baseCurrency)

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

	// Ensure lpIncentive.TotalBlocksPerYear or lpIncentive.EpochNumBlocks are not zero to avoid division by zero
	if lpIncentive.TotalBlocksPerYear.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid inflationary params")
	}

	// Calculate eden amount per epoch
	epochLpsEdenAmount := lpIncentive.EdenAmountPerYear.
		Quo(lpIncentive.TotalBlocksPerYear)

	// Maximum eden based per distribution epoch on maximum APR - 30% by default
	// Allocated for staking per day = (0.3/365)* (total weighted proxy TVL)
	edenDenomPrice := k.GetEdenDenomPrice(ctx, baseCurrency)

	// Ensure edenDenomPrice is not zero to avoid division by zero
	if edenDenomPrice.IsZero() {
		return errorsmod.Wrap(types.ErrNoInflationaryParams, "invalid eden price")
	}

	// Add dexRevenue amount that was tracked by Lp tracker
	dexRevenueLPsAmtPerDistribution = dexRevenueLPsAmtPerDistribution.Add(params.DexRewardsLps.AmountCollectedByOtherTracker)
	// Increase block number
	params.DexRewardsLps.NumBlocks = sdk.NewInt(1)
	// Incrase total dex rewards given
	params.DexRewardsLps.Amount = dexRevenueLPsAmtPerDistribution.Add(gasFeesLPsAmtPerDistribution)
	// Reset amount from other tracker
	params.DexRewardsLps.AmountCollectedByOtherTracker = sdk.ZeroDec()
	// Don't increase Lps rewards blocks, it will be increased whenever LP distribution epoch happens.
	params.DexRewardsStakers.AmountCollectedByOtherTracker = params.DexRewardsStakers.AmountCollectedByOtherTracker.Add(dexRevenueStakersAmtPerDistribution)
	k.SetParams(ctx, params)

	totalEdenGivenLP := sdk.ZeroInt()
	totalRewardsGivenLP := sdk.ZeroInt()
	// Process to increase uncomitted token amount of Eden & Eden boost
	k.cmk.IterateCommitments(
		ctx, func(commitments ctypes.Commitments) bool {
			// Commitment owner
			creator := commitments.Creator
			_, err := sdk.AccAddressFromBech32(creator)
			if err != nil {
				// This could be validator address
				return false
			}

			rewardsByUSDCDeposit := sdk.NewCoins()

			newSumEdenRewardsUnClaimed := sdk.ZeroInt()
			newSumDexRewardsUnClaimed := sdk.ZeroInt()

			// Calculate new unclaimed Eden tokens from LpTokens committed, Dex rewards distribution
			// Distribute gas fees to LPs
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			newUnclaimedEdenTokensLp, dexRewardsLp := k.CalcRewardsForLPs(
				ctx, lpIncentive.TotalBlocksPerYear, totalProxyTVL,
				commitments, epochLpsEdenAmount, gasFeesLPsAmtPerDistribution,
			)
			// Total
			totalEdenGivenLP = totalEdenGivenLP.Add(newUnclaimedEdenTokensLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsLp)

			// Sum for each loop
			newSumEdenRewardsUnClaimed = newSumEdenRewardsUnClaimed.Add(newUnclaimedEdenTokensLp)
			newSumDexRewardsUnClaimed = newSumDexRewardsUnClaimed.Add(dexRewardsLp)
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// Calculate new unclaimed Eden tokens from stable stake LpTokens committed, Dex rewards distribution
			// Distribute gas fees to LPs
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			newUnclaimedEdenTokensStableLp, dexRewardsStableLp := k.CalcRewardsForStableStakeLPs(
				ctx, lpIncentive.TotalBlocksPerYear, totalProxyTVL,
				commitments, epochLpsEdenAmount, gasFeesLPsAmtPerDistribution, baseCurrency,
			)

			// Total
			totalEdenGivenLP = totalEdenGivenLP.Add(newUnclaimedEdenTokensStableLp)
			totalRewardsGivenLP = totalRewardsGivenLP.Add(dexRewardsStableLp)

			// Sum for each loop
			newSumEdenRewardsUnClaimed = newSumEdenRewardsUnClaimed.Add(newUnclaimedEdenTokensStableLp)
			newSumDexRewardsUnClaimed = newSumDexRewardsUnClaimed.Add(dexRewardsStableLp)

			// Sub bucket
			rewardsByUSDCDeposit = rewardsByUSDCDeposit.Add(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokensStableLp))
			rewardsByUSDCDeposit = rewardsByUSDCDeposit.Add(sdk.NewCoin(baseCurrency, dexRewardsStableLp))
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			return false
		},
	)

	// Calcualte the remainings
	edenRemainedLP := epochLpsEdenAmount.Sub(totalEdenGivenLP)
	dexRewardsRemainedLP := dexRevenueLPsAmtPerDistribution.Add(gasFeesLPsAmtPerDistribution).Sub(sdk.NewDecFromInt(totalRewardsGivenLP))

	// Fund community the remain coins
	// ----------------------------------
	edenRemainedCoin := sdk.NewDecCoin(ptypes.Eden, edenRemainedLP)
	dexRewardsRemainedCoin := sdk.NewDecCoinFromDec(baseCurrency, dexRewardsRemainedLP)

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(edenRemainedCoin)
	feePool.CommunityPool = feePool.CommunityPool.Add(dexRewardsRemainedCoin)
	k.SetFeePool(ctx, feePool)
	// ----------------------------------

	// Update APR for amm pools
	k.UpdateAmmPoolAPR(ctx, lpIncentive, totalProxyTVL, edenDenomPrice)

	return nil
}
