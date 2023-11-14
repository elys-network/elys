package keeper

import (
	"fmt"

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

// GetCommitments returns a commitments from its index
func (k Keeper) GetCommitments(ctx sdk.Context, creator string) types.Commitments {
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

func (k Keeper) DeductClaimed(ctx sdk.Context, creator string, denom string, amount sdk.Int) (types.Commitments, error) {
	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)

	fmt.Println("commitments", commitments.Claimed.String(), denom, amount.String())
	// Subtract the amount from the claimed balance
	err := commitments.SubClaimed(sdk.NewCoin(denom, amount))
	if err != nil {
		return types.Commitments{}, err
	}
	return commitments, nil
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

func (k Keeper) HandleWithdrawFromCommitment(ctx sdk.Context, commitments *types.Commitments, addr sdk.AccAddress, amount sdk.Coins) error {
	edenAmount := amount.AmountOf(ptypes.Eden)
	edenBAmount := amount.AmountOf(ptypes.EdenB)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, edenAmount))
	commitments.AddClaimed(sdk.NewCoin(ptypes.EdenB, edenBAmount))
	k.SetCommitments(ctx, *commitments)

	withdrawCoins := amount.
		Sub(sdk.NewCoin(ptypes.Eden, edenAmount)).
		Sub(sdk.NewCoin(ptypes.EdenB, edenBAmount))

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, addr.String(), withdrawCoins)

	// Send the coins to the user's account
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
	return err
}

// Withdraw validator's commission to self delegator
func (k Keeper) ProcessWithdrawValidatorCommission(ctx sdk.Context, delegator string, creator string, denom string, amount sdk.Int) error {
	assetProfile, found := k.apKeeper.GetEntry(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	commitments, err := k.DeductCommitments(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Withdraw to the delegated wallet
	addr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	commitments = k.GetCommitments(ctx, delegator)
	err = k.HandleWithdrawFromCommitment(ctx, &commitments, addr, withdrawCoins)
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

// Withdraw Token - USDC
// Only withraw USDC from dexRevenue wallet
func (k Keeper) ProcessWithdrawUSDC(ctx sdk.Context, creator string, denom string, amount sdk.Int) error {
	if denom != ptypes.BaseCurrency {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	assetProfile, found := k.apKeeper.GetEntry(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	commitments, err := k.DeductCommitments(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})

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
