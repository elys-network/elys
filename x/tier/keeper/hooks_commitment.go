package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

// Process commitmentChanged hook
func (k Keeper) CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coins) error {
	k.RetreiveAllPortfolio(ctx, creator)
	return nil
}

// Process eden uncommitted hook
func (k Keeper) EdenUncommitted(ctx sdk.Context, creator string, amount sdk.Coin) error {
	k.RetreiveAllPortfolio(ctx, creator)
	return nil
}

func (k Keeper) BeforeEdenInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}

func (k Keeper) BeforeEdenBInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}

func (k Keeper) BeforeEdenCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}

func (k Keeper) BeforeEdenBCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentive keeper
type CommitmentHooks struct {
	k Keeper
}

var _ commitmenttypes.CommitmentHooks = CommitmentHooks{}

// Return the wrapper struct
func (k Keeper) CommitmentHooks() CommitmentHooks {
	return CommitmentHooks{k}
}

// CommitmentChanged implements CommentmentHook
func (h CommitmentHooks) CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coins) error {
	return h.k.CommitmentChanged(ctx, creator, amount)
}

// EdenUncommitted implements EdenUncommitted
func (h CommitmentHooks) EdenUncommitted(ctx sdk.Context, creator string, amount sdk.Coin) error {
	return h.k.EdenUncommitted(ctx, creator, amount)
}

func (h CommitmentHooks) BeforeEdenInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	return h.k.BeforeEdenInitialCommit(ctx, addr)
}

func (h CommitmentHooks) BeforeEdenBInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	return h.k.BeforeEdenBInitialCommit(ctx, addr)
}

func (h CommitmentHooks) BeforeEdenCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	return h.k.BeforeEdenCommitChange(ctx, addr)
}

func (h CommitmentHooks) BeforeEdenBCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	return h.k.BeforeEdenBCommitChange(ctx, addr)
}
