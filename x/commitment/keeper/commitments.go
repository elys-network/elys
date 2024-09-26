package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// SetCommitments set a specific commitments in the store from its index
func (k Keeper) SetCommitments(ctx sdk.Context, commitments types.Commitments) {
	if !k.HasCommitments(ctx, commitments.GetCreatorAccount()) {
		params := k.GetParams(ctx)
		params.NumberOfCommitments++
		k.SetParams(ctx, params)
	}
	store := ctx.KVStore(k.storeKey)
	key := types.GetCommitmentsKey(commitments.GetCreatorAccount())
	b := k.cdc.MustMarshal(&commitments)
	store.Set(key, b)
}

// GetAllCommitments returns all commitments
func (k Keeper) GetAllCommitments(ctx sdk.Context) (list []*types.Commitments) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

// remove after migration
func (k Keeper) GetAllLegacyCommitments(ctx sdk.Context) (list []*types.Commitments) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyCommitmentsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

// GetCommitments returns a commitments from its index
func (k Keeper) GetCommitments(ctx sdk.Context, creator sdk.AccAddress) types.Commitments {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.GetCommitmentsKey(creator))
	if b == nil {
		return types.Commitments{
			Creator:         creator.String(),
			CommittedTokens: []*types.CommittedTokens{},
			Claimed:         sdk.Coins{},
			VestingTokens:   []*types.VestingTokens{},
		}
	}

	val := types.Commitments{}
	k.cdc.MustUnmarshal(b, &val)
	return val
}

func (k Keeper) HasCommitments(ctx sdk.Context, creator sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCommitmentsKey(creator)
	return store.Has(key)
}

// remove after migration
func (k Keeper) HasLegacyCommitments(ctx sdk.Context, creator string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyCommitmentsKeyPrefix))
	b := store.Get(types.LegacyCommitmentsKey(creator))
	return b != nil
}

// RemoveCommitments removes a commitments from the store
func (k Keeper) RemoveCommitments(ctx sdk.Context, creator sdk.AccAddress) {
	if k.HasCommitments(ctx, creator) {
		params := k.GetParams(ctx)
		params.NumberOfCommitments--
		k.SetParams(ctx, params)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetCommitmentsKey(creator))
}

// remove after migration
func (k Keeper) DeleteLegacyCommitments(ctx sdk.Context, creator string) {
	if k.HasLegacyCommitments(ctx, creator) {
		params := k.GetParams(ctx)
		params.NumberOfCommitments--
		k.SetParams(ctx, params)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyCommitmentsKeyPrefix))
	store.Delete(types.LegacyCommitmentsKey(creator))
}

// IterateCommitments iterates over all Commitments and performs a
// callback.
func (k Keeper) IterateCommitments(ctx sdk.Context, handlerFn func(commitments types.Commitments) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

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
	params := k.GetParams(ctx)
	return int64(params.NumberOfCommitments)
}

func (k Keeper) DeductClaimed(ctx sdk.Context, creator sdk.AccAddress, denom string, amount math.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	// Subtract the amount from the claimed balance
	err := commitments.SubClaimed(sdk.NewCoin(denom, amount))
	if err != nil {
		return types.Commitments{}, err
	}
	return commitments, nil
}

func (k Keeper) BurnEdenBoost(ctx sdk.Context, creator sdk.AccAddress, denom string, amount math.Int) error {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	// if deduction amount is zero
	if amount.IsZero() {
		return nil
	}

	// Subtract the amount from the claimed balance
	claimed := commitments.GetClaimedForDenom(denom)
	claimedRemovalAmount := amount
	if claimed.LT(claimedRemovalAmount) {
		claimedRemovalAmount = claimed
	}
	err := commitments.SubClaimed(sdk.NewCoin(denom, claimedRemovalAmount))
	if err != nil {
		return err
	}

	amount = amount.Sub(claimedRemovalAmount)
	if amount.IsZero() {
		return nil
	}

	committedAmount := commitments.GetCommittedAmountForDenom(denom)
	if committedAmount.LT(amount) {
		amount = committedAmount
	}
	if amount.IsZero() {
		return nil
	}

	err = k.hooks.BeforeEdenBCommitChange(ctx, creator)
	if err != nil {
		return err
	}

	// Subtract the amount from the committed balance
	err = commitments.DeductFromCommitted(denom, amount, uint64(ctx.BlockTime().Unix()), false)
	if err != nil {
		return err
	}

	k.SetCommitments(ctx, commitments)

	err = k.hooks.CommitmentChanged(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})
	if err != nil {
		return err
	}
	return nil
}
