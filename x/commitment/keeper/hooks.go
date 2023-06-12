package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

var _ types.CommitmentHooks = MultiCommitmentHooks{}

// combine multiple commitment hooks, all hook functions are run in array sequence
type MultiCommitmentHooks []types.CommitmentHooks

func NewMultiEpochHooks(hooks ...types.CommitmentHooks) MultiCommitmentHooks {
	return hooks
}

// Committed is called when staker committed his token
func (mh MultiCommitmentHooks) CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coin) {
	for i := range mh {
		mh[i].CommitmentChanged(ctx, creator, amount)
	}
}

// Committed executes the indicated for committed hook
func (k Keeper) AfterCommitmentChange(ctx sdk.Context, creator string, amount sdk.Coin) {
	if k.hooks == nil {
		return
	}
	k.hooks.CommitmentChanged(ctx, creator, amount)
}
