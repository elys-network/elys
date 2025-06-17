package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/commitment/types"
)

var _ types.CommitmentHooks = MultiCommitmentHooks{}

// combine multiple commitment hooks, all hook functions are run in array sequence
type MultiCommitmentHooks []types.CommitmentHooks

func NewMultiCommitmentHooks(hooks ...types.CommitmentHooks) MultiCommitmentHooks {
	return hooks
}

// Committed is called when staker committed his token
func (mh MultiCommitmentHooks) CommitmentChanged(ctx sdk.Context, creator sdk.AccAddress, amount sdk.Coins) error {
	for i := range mh {
		err := mh[i].CommitmentChanged(ctx, creator, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

// Committed is called when staker committed his token
func (mh MultiCommitmentHooks) EdenUncommitted(ctx sdk.Context, creator sdk.AccAddress, amount sdk.Coin) error {
	for i := range mh {
		err := mh[i].EdenUncommitted(ctx, creator, amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mh MultiCommitmentHooks) BeforeEdenInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	for i := range mh {
		err := mh[i].BeforeEdenInitialCommit(ctx, addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mh MultiCommitmentHooks) BeforeEdenBInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	for i := range mh {
		err := mh[i].BeforeEdenBInitialCommit(ctx, addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mh MultiCommitmentHooks) BeforeEdenCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	for i := range mh {
		err := mh[i].BeforeEdenCommitChange(ctx, addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mh MultiCommitmentHooks) BeforeEdenBCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	for i := range mh {
		err := mh[i].BeforeEdenBCommitChange(ctx, addr)
		if err != nil {
			return err
		}
	}
	return nil
}

// Committed executes the indicated for committed hook
func (k Keeper) CommitmentChanged(ctx sdk.Context, creator sdk.AccAddress, amount sdk.Coins) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.CommitmentChanged(ctx, creator, amount)
}

// Committed executes the indicated for committed hook
func (k Keeper) EdenUncommitted(ctx sdk.Context, creator sdk.AccAddress, amount sdk.Coin) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.EdenUncommitted(ctx, creator, amount)
}
