package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Calculate new Eden token amounts based on the given conditions and user's current unclaimed token balance
func (k Keeper) CalcRewardsForStakersByElysStaked(ctx sdk.Context, delAmount math.Int, edenAmountPerDistribution math.Int, dexRevenueAmtForStakersPerDistribution sdk.Dec) (math.Int, math.Int, sdk.Dec) {
	// -----------Eden calculation ---------------------
	// --------------------------------------------------------------
	stakeShare := k.CalcTotalShareOfStaking(delAmount)

	// Calculate newly creating eden amount by its share
	newEdenAllocated := stakeShare.MulInt(edenAmountPerDistribution)

	// --------------------DEX rewards calculation --------------------
	// ----------------------------------------------------------------
	// Calculate dex rewards
	dexRewards := stakeShare.Mul(dexRevenueAmtForStakersPerDistribution).TruncateInt()

	// Calculate only elys staking share
	stakeShareByStakeOnly := k.CalcTotalShareOfStaking(delAmount)
	dexRewardsByStakeOnly := stakeShareByStakeOnly.Mul(dexRevenueAmtForStakersPerDistribution)

	return newEdenAllocated.TruncateInt(), dexRewards, dexRewardsByStakeOnly
}

// Calculate new Eden token amounts based on the given conditions and user's current unclaimed token balance
func (k Keeper) CalcRewardsForStakersByCommitted(ctx sdk.Context, amount math.Int, edenAmountPerEpoch math.Int, dexRevenueAmtForStakers sdk.Dec) (math.Int, math.Int) {
	// -----------Eden calculation ---------------------
	// --------------------------------------------------------------
	stakeShare := k.CalcTotalShareOfStaking(amount)

	// Calculate newly creating eden amount by its share
	newEdenAllocated := stakeShare.MulInt(edenAmountPerEpoch).TruncateInt()

	// --------------------DEX rewards calculation --------------------
	// ----------------------------------------------------------------
	// Calculate dex rewards
	dexRewards := stakeShare.Mul(dexRevenueAmtForStakers).TruncateInt()

	return newEdenAllocated, dexRewards
}

// Calculate new Eden-Boost token amounts based on the given conditions and user's current unclaimed token balance
func (k Keeper) CalculateEdenBoostRewards(
	ctx sdk.Context,
	delAmount math.Int,
	commitments ctypes.Commitments,
	incentiveInfo types.IncentiveInfo,
	edenBoostAPR sdk.Dec,
) (math.Int, math.Int, math.Int) {
	// Get eden commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	params := k.GetParams(ctx)

	// Calculate the portion of each program contribution
	newEdenBByElysStaked := sdk.NewDecFromInt(delAmount).
		Mul(edenBoostAPR).
		MulInt64(params.DistributionInterval).
		QuoInt(incentiveInfo.TotalBlocksPerYear).
		RoundInt()

	newEdenBByEdenCommitted := sdk.NewDecFromInt(edenCommitted).
		Mul(edenBoostAPR).
		MulInt64(params.DistributionInterval).
		QuoInt(incentiveInfo.TotalBlocksPerYear).
		RoundInt()

	newEdenBoost := newEdenBByElysStaked.Add(newEdenBByEdenCommitted)
	return newEdenBoost, newEdenBByElysStaked, newEdenBByEdenCommitted
}
