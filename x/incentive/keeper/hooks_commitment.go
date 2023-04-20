package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
)

// Process commitmentChanged hook
func (k Keeper) CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coin) {
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
func (h CommitmentHooks) CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coin) {
	h.k.CommitmentChanged(ctx, creator, amount)
}
