package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// CommitmentHooks event hooks for commitment processing
type CommitmentHooks interface {
	CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coins) error
	EdenUncommitted(ctx sdk.Context, creator string, amount sdk.Coin) error
	// Hooks for estaking specific
	BeforeEdenInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error
	BeforeEdenBInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error
	BeforeEdenCommitChange(ctx sdk.Context, addr sdk.AccAddress) error
	BeforeEdenBCommitChange(ctx sdk.Context, addr sdk.AccAddress) error
}
