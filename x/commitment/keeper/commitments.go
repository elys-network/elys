package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
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
			Creator:                 creator,
			CommittedTokens:         []*types.CommittedTokens{},
			RewardsUnclaimed:        sdk.Coins{},
			Claimed:                 sdk.Coins{},
			VestingTokens:           []*types.VestingTokens{},
			RewardsByElysUnclaimed:  sdk.Coins{},
			RewardsByEdenUnclaimed:  sdk.Coins{},
			RewardsByEdenbUnclaimed: sdk.Coins{},
			RewardsByUsdcUnclaimed:  sdk.Coins{},
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

func (k Keeper) DeductUnclaimed(ctx sdk.Context, creator string, denom string, amount math.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	// if deduction amount is zero
	if amount.Equal(sdk.ZeroInt()) {
		return commitments, nil
	}

	// Subtract the withdrawn amount from the unclaimed balance
	err := commitments.SubRewardsUnclaimed(sdk.NewCoin(denom, amount))
	if err != nil {
		return types.Commitments{}, err
	}
	return commitments, nil
}

func (k Keeper) BurnEdenBoost(ctx sdk.Context, creator string, denom string, amount math.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	// if deduction amount is zero
	if amount.Equal(sdk.ZeroInt()) {
		return commitments, nil
	}

	// Subtract the amount from the unclaimed balance
	unclaimedRemovalAmount := amount
	rewardUnclaimed := commitments.GetRewardUnclaimedForDenom(denom)
	if rewardUnclaimed.LT(unclaimedRemovalAmount) {
		unclaimedRemovalAmount = rewardUnclaimed
	}
	err := commitments.SubRewardsUnclaimed(sdk.NewCoin(denom, unclaimedRemovalAmount))
	if err != nil {
		return types.Commitments{}, err
	}

	amount = amount.Sub(unclaimedRemovalAmount)
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

	err = commitments.DeductFromCommitted(denom, amount, uint64(ctx.BlockTime().Unix()))
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
	k.AfterCommitmentChange(ctx, commitments.Creator, amount)

	withdrawCoins := amount.
		Sub(sdk.NewCoin(ptypes.Eden, edenAmount)).
		Sub(sdk.NewCoin(ptypes.EdenB, edenBAmount))

	if sendCoins && !withdrawCoins.Empty() {
		return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
	}
	return nil
}

// Update commitments for validator's commission withdrawal to self delegator
func (k Keeper) RecordWithdrawValidatorCommission(ctx sdk.Context, delegator string, creator string, denom string, amount math.Int) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return errorsmod.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	commitments, err := k.DeductUnclaimed(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Withdraw to the delegated wallet
	addr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	commitments = k.GetCommitments(ctx, delegator)
	err = k.HandleWithdrawFromCommitment(ctx, &commitments, withdrawCoins, false, addr)
	if err != nil {
		return err
	}

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

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
		k.AfterCommitmentChange(ctx, delegator, sdk.Coins{})

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
