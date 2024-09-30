package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Creating a commitment object for a delegator if one does not exist:
func (k Keeper) BeforeDelegationCreated(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Updating portfolio on delegation changes
func (k Keeper) AfterDelegationModified(goCtx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.RetrieveAllPortfolio(ctx, delAddr)
	return nil
}

func (k Keeper) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (k Keeper) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// ________________________________________________________________________________________

// StakingHooks wrapper struct for slashing keeper
type StakingHooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = StakingHooks{}

// Return the wrapper struct
func (k Keeper) StakingHooks() StakingHooks {
	return StakingHooks{k}
}

// staking StakingHooks
// Must be called when a validator is created
func (h StakingHooks) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}

// staking StakingHooks
// Must be called when a validator is created
func (h StakingHooks) AfterUnbondingInitiated(ctx context.Context, id uint64) error {
	return nil
}

// Must be called when a validator's state changes
func (h StakingHooks) BeforeValidatorModified(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a validator is deleted
func (h StakingHooks) AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a validator is bonded
func (h StakingHooks) AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a validator begins unbonding
func (h StakingHooks) AfterValidatorBeginUnbonding(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a delegation is created
func (h StakingHooks) BeforeDelegationCreated(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.BeforeDelegationCreated(ctx, delAddr, valAddr)
}

// Must be called when a delegation's shares are modified
func (h StakingHooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}

// Must be called when a delegation is removed
func (h StakingHooks) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}

func (h StakingHooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.AfterDelegationModified(ctx, delAddr, valAddr)
}

func (h StakingHooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction sdkmath.LegacyDec) error {
	return nil
}
