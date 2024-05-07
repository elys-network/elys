package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// SetCommitments set a specific commitments in the store from its index
func (k Keeper) SetCommitments(ctx sdk.Context, commitments types.Commitments) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))
	b := k.cdc.MustMarshal(&commitments)
	store.Set(types.CommitmentsKey(commitments.Creator), b)
}

// GetAllCommitments returns all commitments
func (k Keeper) GetAllCommitments(ctx sdk.Context) (list []*types.Commitments) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

// GetAllLegacyCommitments returns all legacy commitments
func (k Keeper) GetAllLegacyCommitments(ctx sdk.Context) (list []*types.LegacyCommitments) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyCommitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

// GetCommitments returns a commitments from its index
func (k Keeper) GetCommitments(ctx sdk.Context, creator string) types.Commitments {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))

	b := store.Get(types.CommitmentsKey(creator))
	if b == nil {
		return types.Commitments{
			Creator:         creator,
			CommittedTokens: []*types.CommittedTokens{},
			Claimed:         sdk.Coins{},
			VestingTokens:   []*types.VestingTokens{},
		}
	}

	val := types.Commitments{}
	k.cdc.MustUnmarshal(b, &val)
	return val
}

// RemoveCommitments removes a commitments from the store
func (k Keeper) RemoveCommitments(ctx sdk.Context, creator string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))
	store.Delete(types.CommitmentsKey(creator))
}

// IterateCommitments iterates over all Commitments and performs a
// callback.
func (k Keeper) IterateCommitments(ctx sdk.Context, handlerFn func(commitments types.Commitments) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))

	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var commitments types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &commitments)

		if handlerFn(commitments) {
			break
		}
	}
}

// NumberOfCommitments returns total number of commitment items
func (k Keeper) TotalNumberOfCommitments(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))

	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	numberOfCommitments := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		numberOfCommitments++
	}
	return numberOfCommitments
}

func (k Keeper) DeductClaimed(ctx sdk.Context, creator string, denom string, amount math.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	// Subtract the amount from the claimed balance
	err := commitments.SubClaimed(sdk.NewCoin(denom, amount))
	if err != nil {
		return types.Commitments{}, err
	}
	return commitments, nil
}

func (k Keeper) BurnEdenBoost(ctx sdk.Context, creator string, denom string, amount math.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	addr := sdk.MustAccAddressFromBech32(creator)
	err := k.hooks.BeforeEdenBCommitChange(ctx, addr)
	if err != nil {
		return commitments, err
	}

	// if deduction amount is zero
	if amount.Equal(sdk.ZeroInt()) {
		return commitments, nil
	}

	// Subtract the amount from the claimed balance
	claimed := commitments.GetClaimedForDenom(denom)
	claimedRemovalAmount := amount
	if claimed.LT(claimedRemovalAmount) {
		claimedRemovalAmount = claimed
	}
	err = commitments.SubClaimed(sdk.NewCoin(denom, claimedRemovalAmount))
	if err != nil {
		return types.Commitments{}, err
	}

	amount = amount.Sub(claimedRemovalAmount)
	if amount.Equal(sdk.ZeroInt()) {
		return commitments, nil
	}

	committedAmount := commitments.GetCommittedAmountForDenom(denom)
	if committedAmount.LT(amount) {
		amount = committedAmount
	}

	// Subtract the amount from the committed balance
	err = commitments.DeductFromCommitted(denom, amount, uint64(ctx.BlockTime().Unix()))
	if err != nil {
		return types.Commitments{}, err
	}

	err = k.hooks.CommitmentChanged(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})
	if err != nil {
		return types.Commitments{}, err
	}
	return commitments, nil
}

func (k Keeper) HandleWithdrawFromCommitment(ctx sdk.Context, commitments *types.Commitments, amount sdk.Coins, sendCoins bool, addr sdk.AccAddress) error {
	edenAmount := amount.AmountOf(ptypes.Eden)
	edenBAmount := amount.AmountOf(ptypes.EdenB)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, edenAmount))
	commitments.AddClaimed(sdk.NewCoin(ptypes.EdenB, edenBAmount))
	k.SetCommitments(ctx, *commitments)

	// Emit Hook commitment changed
	err := k.CommitmentChanged(ctx, commitments.Creator, amount)
	if err != nil {
		return err
	}

	withdrawCoins := amount.
		Sub(sdk.NewCoin(ptypes.Eden, edenAmount)).
		Sub(sdk.NewCoin(ptypes.EdenB, edenBAmount))

	if sendCoins && !withdrawCoins.Empty() {
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
	}
	return nil
}

// Process delegation hook - create commitment entities for delegator and validator
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delegator string, validator string) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	_, err = sdk.ValAddressFromBech32(validator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert validator address from bech32")
	}

	/***********************************************************/
	////////////////// Delegator entity //////////////////////////
	/***********************************************************/
	// Get the Commitments for the delegator
	commitments := k.GetCommitments(ctx, delegator)
	if commitments.IsEmpty() {
		k.SetCommitments(ctx, commitments)

		// Emit Hook commitment changed
		err := k.CommitmentChanged(ctx, delegator, sdk.Coins{})
		if err != nil {
			return err
		}

		// Emit blockchain event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCommitmentChanged,
				sdk.NewAttribute(types.AttributeCreator, delegator),
				sdk.NewAttribute(types.AttributeAmount, sdk.ZeroInt().String()),
			),
		)
	}

	/***************************************************************/
	////////////////////// Validator entity /////////////////////////
	// Get the Commitments for the validator
	commitments = k.GetCommitments(ctx, validator)
	if commitments.IsEmpty() {
		k.SetCommitments(ctx, commitments)

		// Emit blockchain event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCommitmentChanged,
				sdk.NewAttribute(types.AttributeCreator, validator),
				sdk.NewAttribute(types.AttributeAmount, sdk.ZeroInt().String()),
			),
		)
	}

	return nil
}
