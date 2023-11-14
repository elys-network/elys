package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
) types.Commitments {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentsKeyPrefix))

	b := store.Get(types.CommitmentsKey(creator))
	if b == nil {
		return types.Commitments{
			Creator:          creator,
			CommittedTokens:  []*types.CommittedTokens{},
			RewardsUnclaimed: sdk.Coins{},
			Claimed:          sdk.Coins{},
			VestingTokens:    []*types.VestingTokens{},
		}
	}

	val := types.Commitments{}
	k.cdc.MustUnmarshal(b, &val)
	return val
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
	commitments := k.GetCommitments(ctx, creator)

	// if deduction amount is zero
	if amount.Equal(sdk.ZeroInt()) {
		return commitments, nil
	}

	// Get user's unclaimed reward
	rewardUnclaimed := commitments.GetRewardUnclaimedForDenom(denom)

	unclaimedRemovalAmount := amount

	// Check if there are enough unclaimed rewards to withdraw
	if rewardUnclaimed.LT(unclaimedRemovalAmount) {
		// Calculate the difference between the requested amount and the available unclaimed balance
		difference := unclaimedRemovalAmount.Sub(rewardUnclaimed)

		err := commitments.DeductFromCommitted(denom, difference, uint64(ctx.BlockTime().Unix()))
		if err != nil {
			return types.Commitments{}, err
		}

		unclaimedRemovalAmount = rewardUnclaimed
	}

	// Subtract the withdrawn amount from the unclaimed balance
	err := commitments.SubRewardsUnclaimed(sdk.NewCoin(denom, unclaimedRemovalAmount))
	if err != nil {
		return types.Commitments{}, err
	}
	return commitments, nil
}
