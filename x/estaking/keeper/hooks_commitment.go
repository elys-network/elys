package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Process commitmentChanged hook
func (k Keeper) CommitmentChanged(ctx sdk.Context, creator string, amount sdk.Coins) error {
	params := k.GetParams(ctx)
	addr := sdk.MustAccAddressFromBech32(creator)

	if !amount.AmountOf(ptypes.Eden).IsZero() {
		edenValAddr, err := sdk.ValAddressFromBech32(params.EdenCommitVal)
		if err != nil {
			panic(err)
		}
		err = k.Keeper.Hooks().AfterDelegationModified(ctx, addr, edenValAddr)
		if err != nil {
			return err
		}
	}

	if !amount.AmountOf(ptypes.EdenB).IsZero() {
		edenBValAddr, err := sdk.ValAddressFromBech32(params.EdenbCommitVal)
		if err != nil {
			panic(err)
		}
		err = k.Keeper.Hooks().AfterDelegationModified(ctx, addr, edenBValAddr)
		if err != nil {
			return err
		}
	}
	return nil

}

// Process eden uncommitted hook
func (k Keeper) EdenUncommitted(ctx sdk.Context, creator string, amount sdk.Coin) error {
	return nil
}

func (k Keeper) BeforeEdenInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	params := k.GetParams(ctx)
	edenValAddr, err := sdk.ValAddressFromBech32(params.EdenCommitVal)
	if err != nil {
		panic(err)
	}
	err = k.Keeper.Hooks().BeforeDelegationCreated(ctx, addr, edenValAddr)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) BeforeEdenBInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	params := k.GetParams(ctx)
	edenBValAddr, err := sdk.ValAddressFromBech32(params.EdenbCommitVal)
	if err != nil {
		panic(err)
	}
	err = k.Keeper.Hooks().BeforeDelegationCreated(ctx, addr, edenBValAddr)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) BeforeEdenCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {

	params := k.GetParams(ctx)
	edenValAddr, err := sdk.ValAddressFromBech32(params.EdenCommitVal)
	if err != nil {
		panic(err)
	}
	err = k.Keeper.Hooks().BeforeDelegationSharesModified(ctx, addr, edenValAddr)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) BeforeEdenBCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	params := k.GetParams(ctx)
	edenBValAddr, err := sdk.ValAddressFromBech32(params.EdenbCommitVal)
	if err != nil {
		panic(err)
	}
	err = k.Keeper.Hooks().BeforeDelegationSharesModified(ctx, addr, edenBValAddr)
	if err != nil {
		return err
	}
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
	return nil
}

func (h CommitmentHooks) BeforeEdenBInitialCommit(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}

func (h CommitmentHooks) BeforeEdenCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}

func (h CommitmentHooks) BeforeEdenBCommitChange(ctx sdk.Context, addr sdk.AccAddress) error {
	return nil
}
