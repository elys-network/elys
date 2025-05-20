package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/elys-network/elys/v4/x/epochs/types"
)

// EpochHooks wrapper struct for incentive keeper
type EpochHooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = EpochHooks{}

// Return the wrapper struct
func (k Keeper) EpochHooks() EpochHooks {
	return EpochHooks{k}
}

// BeforeEpochStart implements EpochHooks
func (h EpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := h.k.GetParams(ctx)
	if epochIdentifier == params.ProviderVestingEpochIdentifier {
		return h.k.ClaimAndVestProviderStakingRewards(ctx)
	}
	return nil
}

// AfterEpochEnd implements EpochHooks
func (h EpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}
