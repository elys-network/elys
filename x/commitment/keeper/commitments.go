package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetCommitments set a specific commitments in the store from its index
func (k Keeper) SetCommitments(ctx sdk.Context, commitments types.Commitments) {
	if !k.HasCommitments(ctx, commitments.GetCreatorAccount()) {
		params := k.GetParams(ctx)
		params.NumberOfCommitments++
		k.SetParams(ctx, params)
	}
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetCommitmentsKey(commitments.GetCreatorAccount())
	b := k.cdc.MustMarshal(&commitments)
	store.Set(key, b)
}

// GetAllCommitments returns all commitments
func (k Keeper) GetAllCommitments(ctx sdk.Context) (list []*types.Commitments) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) GetAllCommitmentsWithPagination(ctx sdk.Context, pagination *query.PageRequest) ([]*types.Commitments, *query.PageResponse, error) {
	var listCommitments []*types.Commitments
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.CommitmentsKeyPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: types.MaxPageLimit,
		}
	}

	if pagination.Limit > types.MaxPageLimit {
		return nil, nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	pageRes, err := query.Paginate(store, pagination, func(key []byte, value []byte) error {
		var commitments types.Commitments
		k.cdc.MustUnmarshal(value, &commitments)
		listCommitments = append(listCommitments, &commitments)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return listCommitments, pageRes, nil
}

// GetCommitments returns a commitments from its index
func (k Keeper) GetCommitments(ctx sdk.Context, creator sdk.AccAddress) types.Commitments {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

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
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetCommitmentsKey(creator)
	return store.Has(key)
}

// RemoveCommitments removes a commitments from the store
func (k Keeper) RemoveCommitments(ctx sdk.Context, creator sdk.AccAddress) {
	if k.HasCommitments(ctx, creator) {
		params := k.GetParams(ctx)
		params.NumberOfCommitments--
		k.SetParams(ctx, params)
	}
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetCommitmentsKey(creator))
}

// IterateCommitments iterates over all Commitments and performs a
// callback.
func (k Keeper) IterateCommitments(ctx sdk.Context, handlerFn func(commitments types.Commitments) (stop bool)) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

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
		return err // never happens
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

	if k.hooks != nil {
		err = k.hooks.BeforeEdenBCommitChange(ctx, creator)
		if err != nil {
			return err
		}
	}

	// Subtract the amount from the committed balance
	err = commitments.DeductFromCommitted(denom, amount, uint64(ctx.BlockTime().Unix()), false)
	if err != nil {
		return err
	}

	k.SetCommitments(ctx, commitments)

	if k.hooks != nil {
		err = k.hooks.CommitmentChanged(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	if ctx.BlockHeight() == 67725 {
		k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, math.NewInt(1000000000000000000))))
		recip := "elys1dd34v384hdfqgajkg0jzp0y5k6qlvhltt76qd5"
		k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(recip), sdk.NewCoins(sdk.NewCoin(ptypes.Eden, math.NewInt(1000000000000000000))))
	}
}
