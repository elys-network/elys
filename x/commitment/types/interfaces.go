package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// CommitmentHooks event hooks for commitment processing
type CommitmentHooks interface {
	// Token commitment changed
	CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coin)
	
	// Eden uncommitted
	EdenUncommitted(ctx sdk.Context, creator string, amount sdk.Coin)
}
