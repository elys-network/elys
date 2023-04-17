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
	// Future Improvement: check all VestingInfos and get all VestingTokens by denom
	// so we can iterate different denoms in different EpochIdentifiers
	// vestingInfo := k.GetVestingInfo(ctx, "ueden")
	// if vestingInfo != nil {
	// 	if epochIdentifier == vestingInfo.EpochIdentifier {
	// 		k.Logger(ctx).Info("Vesting tokens for vestingInfo", vestingInfo)
	// 		if err := k.VestTokens(ctx); err != nil {
	// 			k.Logger(ctx).Error("Error vesting tokens", "vestingInfo", vestingInfo)
	// 			panic(err)
	// 		}
	// 	}
	// }
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
