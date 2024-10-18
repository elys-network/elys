package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

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
func (h StakingHooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}

// staking StakingHooks
// Must be called when a validator is created
func (h StakingHooks) AfterUnbondingInitiated(ctx sdk.Context, id uint64) error {
	return nil
}

// Must be called when a validator's state changes
func (h StakingHooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a validator is deleted
func (h StakingHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a validator is bonded
func (h StakingHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a validator begins unbonding
func (h StakingHooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a delegation is created
func (h StakingHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a delegation's shares are modified
func (h StakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

// Must be called when a delegation is removed
func (h StakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (h StakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	h.k.RetrieveAllPortfolio(ctx, delAddr)
	return nil
}

func (h StakingHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	return nil
}
