package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// EndBlocker of incentive module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// Rewards distribution
	k.ProcessRewardsDistribution(ctx)
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
	stakeIncentive.CurrentEpochInBlocks = stakeIncentive.CurrentEpochInBlocks.Add(sdk.OneInt())
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

// Update unclaimed token amount
// Called back through epoch hook
func (k Keeper) UpdateStakersRewardsUnclaimed(ctx sdk.Context, stakeIncentive types.IncentiveInfo) error {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Calculate each portion of Gas fees collected - stakers, LPs
	totalFeesCollected := sdk.Coins{} // TODO: calculate USDC amount
	gasFeeCollectedDec := sdk.NewDecCoinsFromCoins(totalFeesCollected...)
	rewardPortionForStakers := k.GetParams(ctx).RewardPortionForStakers
	gasFeesForStakers := gasFeeCollectedDec.MulDecTruncate(rewardPortionForStakers)

	// Sum Dex revenue for stakers + Gas fees for stakers and name it dex Revenus for stakers
	// But won't sum dex revenue for LPs and gas fees for LPs as the LP revenue will be rewared by pool.
	dexRevenueForStakersPerDistribution := gasFeesForStakers

	// USDC amount in sdk.Dec type
	dexRevenueStakersAmtPerDistribution := dexRevenueForStakersPerDistribution.AmountOf(baseCurrency)

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

	cparams := k.commKeeper.GetParams(ctx)
	totalElysEdenEdenBStake := k.Keeper.TotalBondedTokens(ctx).
		Add(cparams.TotalCommitted.AmountOf(ptypes.Eden)).
		Add(cparams.TotalCommitted.AmountOf(ptypes.EdenB))

	epochStakersMaxEdenAmount := params.MaxEdenRewardAprStakers.
		MulInt(totalElysEdenEdenBStake).
		QuoInt(stakeIncentive.TotalBlocksPerYear)

	// Use min amount (eden allocation from tokenomics and max apr based eden amount)
	epochStakersEdenAmount = sdk.MinInt(epochStakersEdenAmount, epochStakersMaxEdenAmount.TruncateInt())

	// Set block number and total dex rewards given
	params.DexRewardsStakers.NumBlocks = sdk.OneInt()
	params.DexRewardsStakers.Amount = dexRevenueStakersAmtPerDistribution

	k.SetParams(ctx, params)

	coins := sdk.NewCoins(
		sdk.NewCoin(baseCurrency, dexRevenueStakersAmtPerDistribution.RoundInt()),
		sdk.NewCoin(ptypes.Eden, epochStakersEdenAmount),
	)
	return k.commKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, types.ModuleName, coins.Sort())
}
