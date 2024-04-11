package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

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

	// Calculate the portion of each program contribution
	newEdenBByElysStaked := sdk.NewDecFromInt(delAmount).
		Mul(edenBoostAPR).
		QuoInt(incentiveInfo.TotalBlocksPerYear).
		RoundInt()

	newEdenBByEdenCommitted := sdk.NewDecFromInt(edenCommitted).
		Mul(edenBoostAPR).
		QuoInt(incentiveInfo.TotalBlocksPerYear).
		RoundInt()

	newEdenBoost := newEdenBByElysStaked.Add(newEdenBByEdenCommitted)
	return newEdenBoost, newEdenBByElysStaked, newEdenBByEdenCommitted
}
