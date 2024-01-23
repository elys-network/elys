package keeper

import (
	"errors"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// EndBlocker of incentive module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Elys staked amount tracking
	k.ProcessElysStakedTracking(ctx)

	// Rewards distribution
	k.ProcessRewardsDistribution(ctx)
}

// Elys staked amount tracking
func (k Keeper) ProcessElysStakedTracking(ctx sdk.Context) {
	params := k.GetParams(ctx)
	// Update Elys staked amount every n blocks
	if params.ElysStakeTrackingRate == 0 || ctx.BlockHeight()%params.ElysStakeTrackingRate != 0 {
		return
	}

	// Track the amount of Elys staked
	k.cmk.IterateCommitments(
		ctx, func(commitments ctypes.Commitments) bool {
			// Commitment owner
			creator := commitments.Creator
			_, err := sdk.AccAddressFromBech32(creator)
			if err != nil {
				// This could be validator address
				return false
			}

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			elysStaked := types.ElysStaked{
				Address: creator,
				Amount:  delegatedAmt,
			}

			// Set Elys staked amount
			k.SetElysStaked(ctx, elysStaked)

			return false
		},
	)
}

// Rewards distribution
func (k Keeper) ProcessRewardsDistribution(ctx sdk.Context) {
	// Read tokenomics time based inflation params and update incentive module params.
	if !k.ProcessUpdateIncentiveParams(ctx) {
		ctx.Logger().Error("Invalid tokenomics params", "error", errors.New("Invalid tokenomics params"))
		return
	}

	stakerEpoch, stakeIncentive := k.IsStakerRewardsDistributionEpoch(ctx)
	if stakerEpoch {
		err := k.UpdateStakersRewardsUnclaimed(ctx, stakeIncentive)
		if err != nil {
			ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
		}
	}

	lpsEpoch, lpIncentive := k.IsLPRewardsDistributionEpoch(ctx)
	if lpsEpoch {
		err := k.UpdateLPRewardsUnclaimed(ctx, lpIncentive)
		if err != nil {
			ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
		}
	}
}

func (k Keeper) ProcessUpdateIncentiveParams(ctx sdk.Context) bool {
	listTimeBasedInflations := k.tokenomicsKeeper.GetAllTimeBasedInflation(ctx)
	if len(listTimeBasedInflations) < 1 {
		return false
	}

	params := k.GetParams(ctx)

	// Ensure distribution epoch is not zero to avoid division by zero
	if params.DistributionEpochForLpsInBlocks == 0 || params.DistributionEpochForStakersInBlocks == 0 {
		return false
	}

	for _, inflation := range listTimeBasedInflations {
		if inflation.StartBlockHeight > uint64(ctx.BlockHeight()) || inflation.EndBlockHeight < uint64(ctx.BlockHeight()) {
			continue
		}

		totalBlocksPerYear := sdk.NewInt(int64(inflation.EndBlockHeight - inflation.StartBlockHeight + 1))

		// ------------- LP Incentive parameter -------------
		// ptypes.DaysPerYear is guaranteed to be positive as it is defined as a constant
		EpochNumBlocks := totalBlocksPerYear.Quo(sdk.NewInt(ptypes.DaysPerYear))
		totalDistributionEpochPerYear := totalBlocksPerYear.Quo(sdk.NewInt(params.DistributionEpochForLpsInBlocks))
		// If totalDistributionEpochPerYear is zero, we skip this inflation to avoid division by zero
		if totalBlocksPerYear == sdk.ZeroInt() {
			continue
		}
		currentEpochInBlocks := sdk.NewInt(ctx.BlockHeight() - int64(inflation.StartBlockHeight)).Mul(totalDistributionEpochPerYear).Quo(totalBlocksPerYear)
		maxEdenPerAllocation := sdk.NewInt(int64(inflation.Inflation.LmRewards)).Mul(EpochNumBlocks).Quo(totalBlocksPerYear)
		incentiveInfo := types.IncentiveInfo{
			// reward amount in eden for 1 year
			EdenAmountPerYear: sdk.NewInt(int64(inflation.Inflation.LmRewards)),
			// starting block height of the distribution
			DistributionStartBlock: sdk.NewInt(int64(inflation.StartBlockHeight)),
			// distribution duration - block number per year
			TotalBlocksPerYear: totalBlocksPerYear,
			// we set block numbers in 24 hrs
			EpochNumBlocks: EpochNumBlocks,
			// maximum eden allocation per day that won't exceed 30% apr
			MaxEdenPerAllocation: maxEdenPerAllocation,
			// number of block intervals that distribute rewards.
			DistributionEpochInBlocks: sdk.NewInt(params.DistributionEpochForLpsInBlocks),
			// current epoch in block number
			CurrentEpochInBlocks: currentEpochInBlocks,
			// eden boost apr (0-1) range
			EdenBoostApr: sdk.NewDec(1),
		}

		if len(params.LpIncentives) == 0 {
			params.LpIncentives = append(params.LpIncentives, incentiveInfo)
		} else {
			// If any of block number related parameter changed, we re-calculate the current epoch
			if params.LpIncentives[0].DistributionStartBlock != incentiveInfo.DistributionStartBlock ||
				params.LpIncentives[0].TotalBlocksPerYear != incentiveInfo.TotalBlocksPerYear ||
				params.LpIncentives[0].DistributionEpochInBlocks != incentiveInfo.DistributionEpochInBlocks {
				params.LpIncentives[0].CurrentEpochInBlocks = currentEpochInBlocks
			}
			params.LpIncentives[0].EdenAmountPerYear = incentiveInfo.EdenAmountPerYear
			params.LpIncentives[0].DistributionStartBlock = incentiveInfo.DistributionStartBlock
			params.LpIncentives[0].TotalBlocksPerYear = incentiveInfo.TotalBlocksPerYear
			params.LpIncentives[0].EpochNumBlocks = incentiveInfo.EpochNumBlocks
			params.LpIncentives[0].DistributionEpochInBlocks = incentiveInfo.DistributionEpochInBlocks
			params.LpIncentives[0].EdenBoostApr = incentiveInfo.EdenBoostApr
		}

		// ------------- Stakers parameter -------------
		totalDistributionEpochPerYear = totalBlocksPerYear.Quo(sdk.NewInt(params.DistributionEpochForStakersInBlocks))
		currentEpochInBlocks = sdk.NewInt(ctx.BlockHeight() - int64(inflation.StartBlockHeight)).Mul(totalDistributionEpochPerYear).Quo(totalBlocksPerYear)
		maxEdenPerAllocation = sdk.NewInt(int64(inflation.Inflation.IcsStakingRewards)).Mul(EpochNumBlocks).Quo(totalBlocksPerYear)
		incentiveInfo = types.IncentiveInfo{
			// reward amount in eden for 1 year
			EdenAmountPerYear: sdk.NewInt(int64(inflation.Inflation.IcsStakingRewards)),
			// starting block height of the distribution
			DistributionStartBlock: sdk.NewInt(int64(inflation.StartBlockHeight)),
			// distribution duration - block number per year
			TotalBlocksPerYear: totalBlocksPerYear,
			// we set block numbers in 24 hrs
			EpochNumBlocks: EpochNumBlocks,
			// maximum eden allocation per day that won't exceed 30% apr
			MaxEdenPerAllocation: maxEdenPerAllocation,
			// number of block intervals that distribute rewards.
			DistributionEpochInBlocks: sdk.NewInt(params.DistributionEpochForStakersInBlocks),
			// current epoch in block number
			CurrentEpochInBlocks: currentEpochInBlocks,
			// eden boost apr (0-1) range
			EdenBoostApr: sdk.NewDec(1),
		}

		if len(params.StakeIncentives) == 0 {
			params.StakeIncentives = append(params.StakeIncentives, incentiveInfo)
		} else {
			// If any of block number related parameter changed, we re-calculate the current epoch
			if params.StakeIncentives[0].DistributionStartBlock != incentiveInfo.DistributionStartBlock ||
				params.StakeIncentives[0].TotalBlocksPerYear != incentiveInfo.TotalBlocksPerYear ||
				params.StakeIncentives[0].DistributionEpochInBlocks != incentiveInfo.DistributionEpochInBlocks {
				params.StakeIncentives[0].CurrentEpochInBlocks = currentEpochInBlocks
			}
			params.StakeIncentives[0].EdenAmountPerYear = incentiveInfo.EdenAmountPerYear
			params.StakeIncentives[0].DistributionStartBlock = incentiveInfo.DistributionStartBlock
			params.StakeIncentives[0].TotalBlocksPerYear = incentiveInfo.TotalBlocksPerYear
			params.StakeIncentives[0].EpochNumBlocks = incentiveInfo.EpochNumBlocks
			params.StakeIncentives[0].DistributionEpochInBlocks = incentiveInfo.DistributionEpochInBlocks
			params.StakeIncentives[0].EdenBoostApr = incentiveInfo.EdenBoostApr
		}
		break
	}

	k.SetParams(ctx, params)
	return true
}

func (k Keeper) IsStakerRewardsDistributionEpoch(ctx sdk.Context) (bool, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false, types.IncentiveInfo{}
	}

	// If we don't have enough params
	if len(params.StakeIncentives) < 1 {
		return false, types.IncentiveInfo{}
	}

	// Incentive params initialize
	stakeIncentive := params.StakeIncentives[0]
	if ctx.BlockHeight()%stakeIncentive.DistributionEpochInBlocks.Int64() != 0 {
		return false, types.IncentiveInfo{}
	}

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())
	if stakeIncentive.DistributionStartBlock.GT(curBlockHeight) {
		return false, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.CurrentEpochInBlocks = stakeIncentive.CurrentEpochInBlocks.Add(stakeIncentive.DistributionEpochInBlocks)
	if stakeIncentive.CurrentEpochInBlocks.GTE(stakeIncentive.TotalBlocksPerYear) || curBlockHeight.GT(stakeIncentive.TotalBlocksPerYear.Add(stakeIncentive.DistributionStartBlock)) {
		if len(params.StakeIncentives) > 1 {
			params.StakeIncentives = params.StakeIncentives[1:]
			k.SetParams(ctx, params)
			return false, types.IncentiveInfo{}
		} else {
			params.StakeIncentives = []types.IncentiveInfo(nil)
			k.SetParams(ctx, params)
			return false, types.IncentiveInfo{}
		}
	}

	params.StakeIncentives[0].CurrentEpochInBlocks = stakeIncentive.CurrentEpochInBlocks
	k.SetParams(ctx, params)

	// return found, stake incentive params
	return true, stakeIncentive
}

func (k Keeper) IsLPRewardsDistributionEpoch(ctx sdk.Context) (bool, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false, types.IncentiveInfo{}
	}

	// If we don't have enough params
	if len(params.LpIncentives) < 1 {
		return false, types.IncentiveInfo{}
	}

	// Incentive params initialize
	lpIncentive := params.LpIncentives[0]
	if ctx.BlockHeight()%lpIncentive.DistributionEpochInBlocks.Int64() != 0 {
		return false, types.IncentiveInfo{}
	}

	curBlockHeight := sdk.NewInt(ctx.BlockHeight())
	if lpIncentive.DistributionStartBlock.GT(curBlockHeight) {
		return false, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	lpIncentive.CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks.Add(lpIncentive.DistributionEpochInBlocks)
	if lpIncentive.CurrentEpochInBlocks.GTE(lpIncentive.TotalBlocksPerYear) || curBlockHeight.GT(lpIncentive.TotalBlocksPerYear.Add(lpIncentive.DistributionStartBlock)) {
		if len(params.LpIncentives) > 1 {
			params.LpIncentives = params.LpIncentives[1:]
			k.SetParams(ctx, params)
			return false, types.IncentiveInfo{}
		} else {
			params.LpIncentives = []types.IncentiveInfo(nil)
			k.SetParams(ctx, params)
			return false, types.IncentiveInfo{}
		}
	}

	params.LpIncentives[0].CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks
	k.SetParams(ctx, params)

	// return found, lp incentive params
	return true, lpIncentive
}
