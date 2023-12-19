package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	errorsmod "cosmossdk.io/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		memKey             storetypes.StoreKey
		paramstore         paramtypes.Subspace
		cmk                types.CommitmentKeeper
		stk                types.StakingKeeper
		tci                *types.TotalCommitmentInfo
		authKeeper         types.AccountKeeper
		bankKeeper         types.BankKeeper
		amm                types.AmmKeeper
		oracleKeeper       types.OracleKeeper
		assetProfileKeeper types.AssetProfileKeeper
		epochsKeeper       types.EpochsKeeper
		stableKeeper       types.StableStakeKeeper
		tokenomicsKeeper   types.TokenomicsKeeper

		feeCollectorName    string // name of the FeeCollector ModuleAccount
		dexRevCollectorName string // name of the Dex Revenue ModuleAccount
		authority           string // gov module addresss
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ck types.CommitmentKeeper,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	amm types.AmmKeeper,
	ok types.OracleKeeper,
	ap types.AssetProfileKeeper,
	epochsKeeper types.EpochsKeeper,
	stableKeeper types.StableStakeKeeper,
	tokenomicsKeeper types.TokenomicsKeeper,
	feeCollectorName string,
	dexRevCollectorName string,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		memKey:              memKey,
		paramstore:          ps,
		cmk:                 ck,
		stk:                 sk,
		tci:                 &types.TotalCommitmentInfo{},
		feeCollectorName:    feeCollectorName,
		dexRevCollectorName: dexRevCollectorName,
		authKeeper:          ak,
		bankKeeper:          bk,
		amm:                 amm,
		oracleKeeper:        ok,
		assetProfileKeeper:  ap,
		epochsKeeper:        epochsKeeper,
		stableKeeper:        stableKeeper,
		tokenomicsKeeper:    tokenomicsKeeper,
		authority:           authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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

	// Fund community pool based on the communtiy tax
	dexRevenueRemainedForStakersPerDistribution := k.UpdateCommunityPool(ctx, dexRevenueForStakersPerDistribution)

	// USDC amount in sdk.Dec type
	dexRevenueLPsAmtPerDistribution := dexRevenueForLpsPerDistribution.AmountOf(baseCurrency)
	dexRevenueStakersAmtPerDistribution := dexRevenueRemainedForStakersPerDistribution.AmountOf(baseCurrency)
	gasFeesLPsAmtPerDistribution := gasFeesForLps.AmountOf(baseCurrency)
	edenBoostAPR := stakeIncentive.EdenBoostApr

	// Calculate eden amount per epoch
	params := k.GetParams(ctx)

	// Ensure stakeIncentive.TotalBlocksPerYear or stakeIncentive.AllocationEpochInBlocks are not zero to avoid division by zero
	if stakeIncentive.TotalBlocksPerYear.IsZero() || stakeIncentive.AllocationEpochInBlocks.IsZero() {
		return errorsmod.Wrap(types.ErrNoNonInflationaryParams, "invalid inflationary params")
	}

	// Calculate
	edenAmountPerEpochStakersPerDay := stakeIncentive.EdenAmountPerYear.Mul(stakeIncentive.AllocationEpochInBlocks).Quo(stakeIncentive.TotalBlocksPerYear)

	// Maximum eden based per distribution epoch on maximum APR - 30% by default
	// Allocated for staking per day = (0.3/365)* ( total elys staked + total Eden committed + total Eden boost committed)
	maxEdenAmountPerStakersPerDay := params.MaxEdenRewardAprStakers.MulInt(k.tci.TotalElysBonded.Add(k.tci.TotalEdenEdenBoostCommitted)).MulInt(stakeIncentive.AllocationEpochInBlocks).QuoInt(stakeIncentive.TotalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	edenAmountPerEpochStakersPerDay = sdk.MinInt(edenAmountPerEpochStakersPerDay, maxEdenAmountPerStakersPerDay.TruncateInt())

	// Calculate eden amount per distribution epoch
	edenAmountPerEpochStakersPerDistribution := edenAmountPerEpochStakersPerDay.Mul(stakeIncentive.DistributionEpochInBlocks).Quo(stakeIncentive.AllocationEpochInBlocks)

	// Track the DEX rewards distribution for stakers
	// Add dexRevenue amount that was tracked by Lp tracker
	dexRevenueStakersAmtPerDistribution = dexRevenueStakersAmtPerDistribution.Add(params.DexRewardsStakers.AmountCollectedByOtherTracker)
	// Increase block number
	params.DexRewardsStakers.NumBlocks = params.DexRewardsStakers.NumBlocks.Add(stakeIncentive.DistributionEpochInBlocks)
	// Incrase total dex rewards given
	params.DexRewardsStakers.Amount = params.DexRewardsStakers.Amount.Add(dexRevenueStakersAmtPerDistribution)
	// Reset amount from other tracker
	params.DexRewardsStakers.AmountCollectedByOtherTracker = sdk.ZeroDec()
	// Don't increase Lps rewards blocks, it will be increased whenever LP distribution epoch happens.
	params.DexRewardsLps.AmountCollectedByOtherTracker = params.DexRewardsLps.AmountCollectedByOtherTracker.Add(dexRevenueLPsAmtPerDistribution).Add(gasFeesLPsAmtPerDistribution)
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
			rewardsByUSDCDeposit := sdk.NewCoins()

			newSumEdenRewardsUnClaimed := sdk.ZeroInt()
			newSumEdenBRewardsUnClaimed := sdk.ZeroInt()
			newSumDexRewardsUnClaimed := sdk.ZeroInt()

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			// Calculate new unclaimed Eden tokens from Eden & Eden boost committed, Dex rewards distribution
			// Distribute gas fees to stakers

			// Calculate new unclaimed Eden tokens from Elys staked
			// ----------------------------------------------------------
			newUnclaimedEdenTokens, dexRewards, dexRewardsByStakers := k.CalculateRewardsForStakersByElysStaked(ctx, delegatedAmt, edenAmountPerEpochStakersPerDistribution, dexRevenueStakersAmtPerDistribution)
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
			newUnclaimedEdenTokens, dexRewards = k.CalculateRewardsForStakersByCommitted(ctx, edenCommitted, edenAmountPerEpochStakersPerDistribution, dexRevenueStakersAmtPerDistribution)

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
			newUnclaimedEdenTokens, dexRewards = k.CalculateRewardsForStakersByCommitted(ctx, edenBoostCommitted, edenAmountPerEpochStakersPerDistribution, dexRevenueStakersAmtPerDistribution)

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

			// ----------------------------------------------------------
			// Give commission to validators ( Eden from stakers and Dex rewards from stakers. )
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			edenCommissionGiven, dexRewardsCommissionGiven := k.GiveCommissionToValidators(ctx, creator, delegatedAmt, newEdenFromElysStaking, dexRewardsByStakers, baseCurrency)

			// Minus the commission amount given
			newSumEdenRewardsUnClaimed = newSumEdenRewardsUnClaimed.Sub(edenCommissionGiven)

			// Minus the commission amount given
			newSumDexRewardsUnClaimed = newSumDexRewardsUnClaimed.Sub(dexRewardsCommissionGiven)
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// We should deduct validator commissions from "reward by elys staking sub bucket"
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			newEdenFromElysStaking = newEdenFromElysStaking.Sub(edenCommissionGiven)
			newDexRewardFromElysStaking = newDexRewardFromElysStaking.Sub(dexRewardsCommissionGiven)

			// Add Eden rewards from Elys staking
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.Eden, newEdenFromElysStaking))
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(baseCurrency, newDexRewardFromElysStaking))
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// Calculate new unclaimed Eden-Boost tokens for staker and Eden token holders
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			newUnclaimedEdenBoostTokens, newUnclaimedEdenBoostFromElysStaking, newUnclaimedEdenBoostFromEdenCommited := k.CalculateEdenBoostRewards(ctx, delegatedAmt, commitments, stakeIncentive, edenBoostAPR)
			rewardsByElysStaking = rewardsByElysStaking.Add(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBoostFromElysStaking))
			rewardsByEdenCommitted = rewardsByEdenCommitted.Add(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBoostFromEdenCommited))

			newSumEdenBRewardsUnClaimed = newSumEdenBRewardsUnClaimed.Add(newUnclaimedEdenBoostTokens)
			// ----------------------------------------------------------
			// ----------------------------------------------------------

			// Update Commitments with new unclaimed token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newSumEdenRewardsUnClaimed, newSumEdenBRewardsUnClaimed, newSumDexRewardsUnClaimed, baseCurrency)

			// Update sub buckets commitment with new unclaimed token amounts
			k.UpdateCommitmentsSubBuckets(ctx, creator, &commitments, rewardsByElysStaking, rewardsByEdenCommitted, rewardsByEdenBCommitted, rewardsByUSDCDeposit)
			return false
		},
	)

	// Calcualte the remainings
	edenRemained := edenAmountPerEpochStakersPerDistribution.Sub(totalEdenGiven)
	dexRewardsRemained := dexRevenueStakersAmtPerDistribution.Sub(sdk.NewDecFromInt(totalRewardsGiven))

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

	// Fund community pool based on the communtiy tax
	dexRevenueRemainedForStakersPerDistribution := k.UpdateCommunityPool(ctx, dexRevenueForStakersPerDistribution)

	// USDC amount in sdk.Dec type
	dexRevenueLPsAmtPerDistribution := dexRevenueForLpsPerDistribution.AmountOf(baseCurrency)
	dexRevenueStakersAmtPerDistribution := dexRevenueRemainedForStakersPerDistribution.AmountOf(baseCurrency)
	gasFeesLPsAmtPerDistribution := gasFeesForLpsPerDistribution.AmountOf(baseCurrency)

	// Proxy TVL
	// Multiplier on each liquidity pool
	// We have 3 pools of 20, 30, 40 TVL
	// We have mulitplier of 0.3, 0.5, 1.0
	// Proxy TVL = 20*0.3+30*0.5+40*1.0
	totalProxyTVL := k.CalculateProxyTVL(ctx, baseCurrency)

	// Ensure lpIncentive.TotalBlocksPerYear or lpIncentive.AllocationEpochInBlocks are not zero to avoid division by zero
	if lpIncentive.TotalBlocksPerYear.IsZero() || lpIncentive.AllocationEpochInBlocks.IsZero() {
		return errorsmod.Wrap(types.ErrNoNonInflationaryParams, "invalid inflationary params")
	}

	// Calculate eden amount per epoch
	edenAmountPerEpochLPsPerDay := lpIncentive.EdenAmountPerYear.Mul(lpIncentive.AllocationEpochInBlocks).Quo(lpIncentive.TotalBlocksPerYear)

	// Track the DEX rewards distribution for stakers
	params := k.GetParams(ctx)

	// Maximum eden based per distribution epoch on maximum APR - 30% by default
	// Allocated for staking per day = (0.3/365)* (total weighted proxy TVL)
	maxEdenAmountPerLpsPerDay := params.MaxEdenRewardAprLps.Mul(totalProxyTVL).MulInt(lpIncentive.AllocationEpochInBlocks).QuoInt(lpIncentive.TotalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	edenAmountPerEpochLPsPerDay = sdk.MinInt(edenAmountPerEpochLPsPerDay, maxEdenAmountPerLpsPerDay.TruncateInt())

	// Calculate Eden amount per distribution epoch
	edenAmountPerEpochLPsPerDistribution := edenAmountPerEpochLPsPerDay.Mul(lpIncentive.DistributionEpochInBlocks).Quo(lpIncentive.AllocationEpochInBlocks)

	// Add dexRevenue amount that was tracked by Lp tracker
	dexRevenueLPsAmtPerDistribution = dexRevenueLPsAmtPerDistribution.Add(params.DexRewardsLps.AmountCollectedByOtherTracker)
	// Increase block number
	params.DexRewardsLps.NumBlocks = params.DexRewardsLps.NumBlocks.Add(lpIncentive.DistributionEpochInBlocks)
	// Incrase total dex rewards given
	params.DexRewardsLps.Amount = params.DexRewardsLps.Amount.Add(dexRevenueLPsAmtPerDistribution).Add(gasFeesLPsAmtPerDistribution)
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

			rewardsByElysStaking := sdk.NewCoins()
			rewardsByEdenCommitted := sdk.NewCoins()
			rewardsByEdenBCommitted := sdk.NewCoins()
			rewardsByUSDCDeposit := sdk.NewCoins()

			newSumEdenRewardsUnClaimed := sdk.ZeroInt()
			newSumEdenBRewardsUnClaimed := sdk.ZeroInt()
			newSumDexRewardsUnClaimed := sdk.ZeroInt()

			// Calculate new unclaimed Eden tokens from LpTokens committed, Dex rewards distribution
			// Distribute gas fees to LPs
			// ----------------------------------------------------------
			// ----------------------------------------------------------
			newUnclaimedEdenTokensLp, dexRewardsLp := k.CalculateRewardsForLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLPsPerDistribution, gasFeesLPsAmtPerDistribution)
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
			newUnclaimedEdenTokensStableLp, dexRewardsStableLp := k.CalculateRewardsForStableStakeLPs(ctx, totalProxyTVL, commitments, edenAmountPerEpochLPsPerDistribution, gasFeesLPsAmtPerDistribution, baseCurrency)

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

			// Update Commitments with new unclaimed token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newSumEdenRewardsUnClaimed, newSumEdenBRewardsUnClaimed, newSumDexRewardsUnClaimed, baseCurrency)

			// Update sub buckets commitment with new unclaimed token amounts
			k.UpdateCommitmentsSubBuckets(ctx, creator, &commitments, rewardsByElysStaking, rewardsByEdenCommitted, rewardsByEdenBCommitted, rewardsByUSDCDeposit)
			return false
		},
	)

	// Calcualte the remainings
	edenRemainedLP := edenAmountPerEpochLPsPerDistribution.Sub(totalEdenGivenLP)
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
	k.UpdateAmmPoolAPR(ctx, lpIncentive, totalProxyTVL)

	return nil
}

// Update commitment record
func (k Keeper) UpdateCommitments(ctx sdk.Context, creator string, commitments *ctypes.Commitments, newUnclaimedEdenTokens sdk.Int, newUnclaimedEdenBoostTokens sdk.Int, dexRewards sdk.Int, baseCurrency string) {
	// Update unclaimed Eden balances in the Commitments structure
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))
	// Update unclaimed Eden-Boost token balances in the Commitments structure
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.EdenB, newUnclaimedEdenBoostTokens))

	// All dex revenue are collected to incentive module in USDC
	// Gas fees (Elys) are also converted into USDC and collected into total dex revenue wallet of incentive module.
	// Update USDC balances in the Commitments structure.
	// These are the rewards from each pool, margin, gas fee.
	commitments.AddRewardsUnclaimed(sdk.NewCoin(baseCurrency, dexRewards))

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update sub bucket commitment record
func (k Keeper) UpdateCommitmentsSubBuckets(ctx sdk.Context, creator string, commitments *ctypes.Commitments, rewardsByElysStaking sdk.Coins, rewardsByEdenCommitted sdk.Coins, rewardsByEdenBCommitted sdk.Coins, rewardsByUSDCDeposit sdk.Coins) {
	// Add to Elys staking bucket
	commitments.AddSubBucketRewardsByElysUnclaimed(rewardsByElysStaking)
	// Add to Eden committed bucket
	commitments.AddSubBucketRewardsByEdenUnclaimed(rewardsByEdenCommitted)
	// Add to EdenB committed bucket
	commitments.AddSubBucketRewardsByEdenBUnclaimed(rewardsByEdenBCommitted)
	// Add to USDC deposit bucket
	commitments.AddSubBucketRewardsByUsdcUnclaimed(rewardsByUSDCDeposit)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Calculate Proxy TVL
func (k Keeper) CalculateProxyTVL(ctx sdk.Context, baseCurrency string) sdk.Dec {
	multipliedShareSum := sdk.ZeroDec()
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, p.GetPoolId())
		if !found {
			return false
		}

		proxyTVL := tvl.Mul(poolInfo.Multiplier)

		// Calculate total pool share by TVL and multiplier
		multipliedShareSum = multipliedShareSum.Add(proxyTVL)

		return false
	})

	//-----------------------------------
	// Handle stable stake pool
	stableStakePoolId := uint64(stabletypes.PoolId)

	// Get pool info from incentive param
	poolInfo, found := k.GetPoolInfo(ctx, stableStakePoolId)
	if !found {
		k.InitStableStakePoolParams(ctx, stableStakePoolId)
		poolInfo, _ = k.GetPoolInfo(ctx, stableStakePoolId)
	}
	tvl := stabletypes.TVL(ctx, k.authKeeper, k.bankKeeper, baseCurrency)
	proxyTVL := sdk.NewDecFromInt(tvl).Mul(poolInfo.Multiplier)
	multipliedShareSum = multipliedShareSum.Add(proxyTVL)

	// return total sum of TVL share using multiplier of all pools
	return multipliedShareSum
}

// Caculate total TVL
func (k Keeper) CalculateTVL(ctx sdk.Context) sdk.Dec {
	TVL := sdk.ZeroDec()

	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}
		TVL = TVL.Add(tvl)
		return false
	})

	return TVL
}

// Update APR for AMM pool
func (k Keeper) UpdateAmmPoolAPR(ctx sdk.Context, lpIncentive types.IncentiveInfo, totalProxyTVL sdk.Dec) {
	// Iterate to calculate total Eden from LpElys, MElys committed
	k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
		tvl, err := p.TVL(ctx, k.oracleKeeper)
		if err != nil {
			return false
		}

		// Get pool Id
		poolId := p.GetPoolId()

		// Get pool info from incentive param
		poolInfo, found := k.GetPoolInfo(ctx, poolId)
		if !found {
			return false
		}

		poolInfo.NumBlocks = poolInfo.NumBlocks.Add(lpIncentive.DistributionEpochInBlocks)
		// Invalid block number
		if poolInfo.NumBlocks.IsZero() {
			return false
		}

		// Calculate Proxy TVL share considering multiplier
		proxyTVL := tvl.Mul(poolInfo.Multiplier)
		poolShare := sdk.ZeroDec()
		if totalProxyTVL.IsPositive() {
			poolShare = proxyTVL.Quo(totalProxyTVL)
		}

		if tvl.IsZero() {
			return false
		}

		// Dex reward Apr per pool =  total accumulated usdc rewards for 7 day * 52/ tvl of pool
		totalLMDexRewardsAllocatedPerWeek := poolInfo.DexRewardAmountGiven.MulInt(lpIncentive.AllocationEpochInBlocks).MulInt(sdk.NewInt(ptypes.DaysPerWeek)).QuoInt(poolInfo.NumBlocks)
		poolInfo.DexApr = totalLMDexRewardsAllocatedPerWeek.MulInt(sdk.NewInt(ptypes.WeeksPerYear)).Quo(tvl)

		// Eden reward Apr per pool = (total LM Eden reward allocated per day*((tvl of pool * multiplier)/total proxy TVL) ) * 365 / TVL of pool
		totalLMEdenRewardsAllocatedPerDay := poolInfo.EdenRewardAmountGiven.Mul(lpIncentive.AllocationEpochInBlocks).Quo(poolInfo.NumBlocks)
		poolInfo.EdenApr = sdk.NewDecFromInt(totalLMEdenRewardsAllocatedPerDay).Mul(poolShare).MulInt(sdk.NewInt(ptypes.DaysPerYear)).Quo(tvl)

		// Update Pool Info
		k.SetPoolInfo(ctx, poolId, poolInfo)

		return false
	})
}

// Get total dex rewards amount from the specified pool
func (k Keeper) GetDexRewardsAmountForPool(ctx sdk.Context, poolId uint64) sdk.Dec {
	if _, found := k.amm.GetPool(ctx, poolId); !found {
		return sdk.ZeroDec()
	}

	lpPortionPercent := k.GetDEXRewardPortionForLPs(ctx)
	if lpPortionPercent.IsZero() {
		return sdk.ZeroDec()
	}

	if k.tci == nil || k.tci.PoolRevenueTrack == nil || len(k.tci.PoolRevenueTrack) < 1 {
		return sdk.ZeroDec()
	}

	// reward tracking key
	trackKey := types.GetPoolRevenueTrackKey(poolId)
	revenue := k.tci.PoolRevenueTrack[trackKey]
	// calculate total dex rewards
	return revenue.Quo(lpPortionPercent)
}
