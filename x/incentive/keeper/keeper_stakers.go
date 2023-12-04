package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Calculate new Eden token amounts based on the given conditions and user's current unclaimed token balance
func (k Keeper) CalculateRewardsForStakersByElysStaked(ctx sdk.Context, delegatedAmt sdk.Int, edenAmountPerDistribution sdk.Int, dexRevenueAmtForStakersPerDistribution sdk.Dec) (sdk.Int, sdk.Int, sdk.Dec) {
	// -----------Eden calculation ---------------------
	// --------------------------------------------------------------
	stakeShare := k.CalculateTotalShareOfStaking(delegatedAmt)

	// Calculate newly creating eden amount by its share
	newEdenAllocated := stakeShare.MulInt(edenAmountPerDistribution)

	// -----------------Fund community Eden token----------------------
	// ----------------------------------------------------------------
	edenCoin := sdk.NewDecCoinFromDec(ptypes.Eden, newEdenAllocated)
	newEdenCoinRemained := k.UpdateCommunityPool(ctx, sdk.DecCoins{edenCoin})

	// Get remained Eden amount
	newEdenAllocated = newEdenCoinRemained.AmountOf(ptypes.Eden)

	// --------------------DEX rewards calculation --------------------
	// ----------------------------------------------------------------
	// Calculate dex rewards
	dexRewards := stakeShare.Mul(dexRevenueAmtForStakersPerDistribution).TruncateInt()

	// Calculate only elys staking share
	stakeShareByStakeOnly := k.CalculateTotalShareOfStaking(delegatedAmt)
	dexRewardsByStakeOnly := stakeShareByStakeOnly.Mul(dexRevenueAmtForStakersPerDistribution)

	return newEdenAllocated.TruncateInt(), dexRewards, dexRewardsByStakeOnly
}

// Calculate new Eden token amounts based on the given conditions and user's current unclaimed token balance
func (k Keeper) CalculateRewardsForStakersByCommitted(ctx sdk.Context, amt sdk.Int, edenAmountPerEpoch sdk.Int, dexRevenueAmtForStakers sdk.Dec) (sdk.Int, sdk.Int) {
	// -----------Eden calculation ---------------------
	// --------------------------------------------------------------
	stakeShare := k.CalculateTotalShareOfStaking(amt)

	// Calculate newly creating eden amount by its share
	newEdenAllocated := stakeShare.MulInt(edenAmountPerEpoch)

	// -----------------Fund community Eden token----------------------
	// ----------------------------------------------------------------
	edenCoin := sdk.NewDecCoinFromDec(ptypes.Eden, newEdenAllocated)
	newEdenCoinRemained := k.UpdateCommunityPool(ctx, sdk.DecCoins{edenCoin})

	// Get remained Eden amount
	newEdenAllocated = newEdenCoinRemained.AmountOf(ptypes.Eden)

	// --------------------DEX rewards calculation --------------------
	// ----------------------------------------------------------------
	// Calculate dex rewards
	dexRewards := stakeShare.Mul(dexRevenueAmtForStakers).TruncateInt()

	return newEdenAllocated.TruncateInt(), dexRewards
}

// Calculate new Eden-Boost token amounts based on the given conditions and user's current unclaimed token balance
func (k Keeper) CalculateEdenBoostRewards(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, incentiveInfo types.IncentiveInfo, edenBoostAPR sdk.Dec) (sdk.Int, sdk.Int, sdk.Int) {
	// Get eden commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Compute eden reward based on above and param factors for each
	totalEden := delegatedAmt.Add(edenCommitted)

	// Ensure incentiveInfo.DistributionEpochInBlocks is not zero to avoid division by zero
	if incentiveInfo.DistributionEpochInBlocks.IsZero() {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
	}

	// Calculate edenBoostAPR % APR for eden boost
	epochNumsPerYear := incentiveInfo.TotalBlocksPerYear.Quo(incentiveInfo.DistributionEpochInBlocks)
	if epochNumsPerYear.IsZero() {
		return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
	}

	apr := edenBoostAPR.TruncateInt()
	newEdenBoost := totalEden.Quo(epochNumsPerYear).Mul(apr)

	// Calculate the portion of each program contribution
	newEdenBoostByElysStaked := sdk.ZeroInt()
	if !totalEden.IsZero() {
		newEdenBoostByElysStaked = sdk.NewDecFromInt(delegatedAmt).QuoInt(totalEden).MulInt(newEdenBoost).TruncateInt()
	}
	newEdenBoostByEdenCommitted := newEdenBoost.Sub(newEdenBoostByElysStaked)

	return newEdenBoost, newEdenBoostByElysStaked, newEdenBoostByEdenCommitted
}
