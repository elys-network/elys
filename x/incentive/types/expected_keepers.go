package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
)

// CommitmentKeeper
type CommitmentKeeper interface {
	// Iterate all commitment
	IterateCommitments(sdk.Context, func(ctypes.Commitments) (stop bool))
	// Initiate commitment according to standard staking
	StandardStakingToken(sdk.Context, string, string, string) error
	// Update commitment
	SetCommitments(sdk.Context, ctypes.Commitments)
	// Get commitment
	GetCommitments(sdk.Context, string) (ctypes.Commitments, bool)
}

// Staking keeper
type StakingKeeper interface {
	TotalBondedTokens(sdk.Context) sdk.Int // total bonded tokens within the validator set
	// iterate through all delegations from one delegator by validator-AccAddress,
	// execute func for each validator
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	// get a particular validator by operator address
	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
	// GetDelegatorDelegations returns a given amount of all the delegations from a delegator.
	GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation)
}
