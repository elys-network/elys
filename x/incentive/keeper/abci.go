// Copyright 2022 Evmos Foundation
// This file is part of the Evmos Network packages.
//
// Evmos is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Evmos packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Evmos packages. If not, see https://github.com/evmos/evmos/blob/main/LICENSE

package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
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
	stakerEpoch, stakeIncentive := k.IsStakerRewardsDistributionEpoch(ctx)
	if stakerEpoch {
		err := k.UpdateStakersRewardsUnclaimed(ctx, stakeIncentive)
		if err != nil {
			ctx.Logger().Error("Failed to update staker rewards unclaimed", "error", err)
		}
	}

	lpsEpoch, lpIncentive := k.IsStakerRewardsDistributionEpoch(ctx)
	if lpsEpoch {
		err := k.UpdateLPRewardsUnclaimed(ctx, lpIncentive)
		if err != nil {
			ctx.Logger().Error("Failed to update lp rewards unclaimed", "error", err)
		}
	}
}

func (k Keeper) IsStakerRewardsDistributionEpoch(ctx sdk.Context) (bool, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false, types.IncentiveInfo{}
	}

	// Update params
	defer k.SetParams(ctx, params)

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
	if stakeIncentive.DistributionStartBlock.LT(curBlockHeight) {
		return false, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.CurrentEpochInBlocks = stakeIncentive.CurrentEpochInBlocks.Add(stakeIncentive.DistributionEpochInBlocks)
	if stakeIncentive.CurrentEpochInBlocks.GTE(stakeIncentive.TotalBlocksPerYear) {
		if len(params.StakeIncentives) > 1 {
			params.StakeIncentives = params.StakeIncentives[1:]
		} else {
			return false, types.IncentiveInfo{}
		}
	}

	// return found, stake, lp incentive params
	return true, stakeIncentive
}

func (k Keeper) IsLPRewardsDistributionEpoch(ctx sdk.Context) (bool, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)
	if ctx.BlockHeight() < 1 {
		return false, types.IncentiveInfo{}
	}

	// Update params
	defer k.SetParams(ctx, params)

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
	if lpIncentive.DistributionStartBlock.LT(curBlockHeight) {
		return false, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	lpIncentive.CurrentEpochInBlocks = lpIncentive.CurrentEpochInBlocks.Add(lpIncentive.DistributionEpochInBlocks)
	if lpIncentive.CurrentEpochInBlocks.GTE(lpIncentive.TotalBlocksPerYear) {
		if len(params.StakeIncentives) > 1 {
			params.LpIncentives = params.LpIncentives[1:]
		} else {
			return false, types.IncentiveInfo{}
		}
	}

	// return found, stake, lp incentive params
	return true, lpIncentive
}
