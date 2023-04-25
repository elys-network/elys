package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
)

// SetCommitments set a specific commitments in the store from its index
func (k Keeper) SetCommitments(ctx sdk.Context, commitments types.Commitments) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))
	b := k.cdc.MustMarshal(&commitments)
	store.Set(types.CommitmentsKey(
		commitments.Creator,
	), b)
}

// GetCommitments returns a commitments from its index
func (k Keeper) GetCommitments(
	ctx sdk.Context,
	creator string,

) (val types.Commitments, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))

	b := store.Get(types.CommitmentsKey(
		creator,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCommitments removes a commitments from the store
func (k Keeper) RemoveCommitments(
	ctx sdk.Context,
	creator string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))
	store.Delete(types.CommitmentsKey(
		creator,
	))
}

// IterateCommitments iterates over all Commitments and performs a
// callback.
func (k Keeper) IterateCommitments(
	ctx sdk.Context, handlerFn func(commitments types.Commitments) (stop bool),
) {
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

func (k Keeper) DeductCommitments(ctx sdk.Context, creator string, denom string, amount sdk.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, creator)
	if !found {
		return types.Commitments{}, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", creator)
	}

	// Get user's uncommitted balance
	uncommittedToken, _ := commitments.GetUncommittedTokensForDenom(denom)

	requestedAmount := amount

	// Check if there are enough uncommitted tokens to withdraw
	if uncommittedToken.Amount.LT(requestedAmount) {
		// Calculate the difference between the requested amount and the available uncommitted balance
		difference := requestedAmount.Sub(uncommittedToken.Amount)

		committedToken, found := commitments.GetCommittedTokensForDenom(denom)
		if found {
			if committedToken.Amount.GTE(difference) {
				// Uncommit the required committed tokens
				committedToken.Amount = committedToken.Amount.Sub(difference)
				requestedAmount = requestedAmount.Sub(difference)
			} else {
				return types.Commitments{}, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "not enough tokens to withdraw")
			}
		}
	}

	// Subtract the withdrawn amount from the uncommitted balance
	uncommittedToken.Amount = uncommittedToken.Amount.Sub(requestedAmount)

	return commitments, nil
}
