package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
)

// BeforeEpochStart performs a no-op
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {

}

// AfterEpochEnd distributes vested tokens at the end of each epoch
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) {
	// Find out incentive param using epochIdentifier and current block timestamp
	foundIncentive, stakeIncentive, lpIncentive := k.GetProperIncentiveParam(ctx, epochIdentifier)

	// If there is no incentive available with the current epoch and timestamp,
	if !foundIncentive {
		return
	}

	// Update unclaimed token amount
	k.UpdateRewardsUnclaimed(ctx, epochIdentifier, stakeIncentive, lpIncentive)
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentive keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeEpochStart implements EpochHooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

// AfterEpochEnd implements EpochHooks
func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
