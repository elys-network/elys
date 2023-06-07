package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Calculate new Eden token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateRewardsForStakers(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, edenAmountPerEpoch sdk.Int, dexRevenueAmtForStakers sdk.Dec) (sdk.Int, sdk.Int, sdk.Dec) {
	// -----------Eden & Eden boost calculation ---------------------
	// --------------------------------------------------------------
	// Get eden commitments and eden boost commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)
	edenBoostCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)

	// compute eden reward based on above and param factors for each
	totalEdenCommittedByStake := delegatedAmt.Add(edenCommitted).Add(edenBoostCommitted)
	stakeShare := k.CalculateTotalShareOfStaking(totalEdenCommittedByStake)

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

	// Calculate only elys staking share
	stakeShareByStakeOnly := k.CalculateTotalShareOfStaking(delegatedAmt)
	dexRewardsByStakeOnly := stakeShareByStakeOnly.Mul(dexRevenueAmtForStakers)

	return newEdenAllocated.TruncateInt(), dexRewards, dexRewardsByStakeOnly
}

// Calculate new Eden-Boost token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateEdenBoostRewards(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, epochIdentifier string, edenBoostAPR int64) sdk.Int {
	// Get eden commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Compute eden reward based on above and param factors for each
	totalEden := delegatedAmt.Add(edenCommitted)

	// Calculate edenBoostAPR % APR for eden boost
	epochNumsPerYear := k.CalculateEpochCountsPerYear(epochIdentifier)

	return totalEden.Quo(sdk.NewInt(epochNumsPerYear)).Quo(sdk.NewInt(100)).Mul(sdk.NewInt(edenBoostAPR))
}
