package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Creating a commitment object for a delegator if one does not exist:
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	// must not run on genesis
	if ctx.BlockHeight() <= 1 {
		return nil
	}

	// Create an entity in commitment module
	k.cmk.StandardStakingToken(ctx, delAddr.String(), valAddr.String(), ptypes.Eden)
	k.cmk.StandardStakingToken(ctx, delAddr.String(), valAddr.String(), ptypes.BaseCurrency)

	return nil
}

// Updating commitments on delegation changes
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return k.BurnEdenBFromElysUnstaking(ctx, delAddr)
}

func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return k.BurnEdenBFromElysUnstaking(ctx, delAddr)
}

func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return k.BurnEdenBFromElysUnstaking(ctx, delAddr)
}

// ________________________________________________________________________________________

// StakingHooks wrapper struct for slashing keeper
type StakingHooks struct {
	k Keeper
}

var (
	_ stakingtypes.StakingHooks = StakingHooks{}
)

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
	return h.k.BeforeDelegationCreated(ctx, delAddr, valAddr)
}

// Must be called when a delegation's shares are modified
func (h StakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}

// Must be called when a delegation is removed
func (h StakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}

func (h StakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.AfterDelegationModified(ctx, delAddr, valAddr)
}

func (h StakingHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	return nil
}
