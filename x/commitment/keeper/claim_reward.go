package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// Update commitments for claim reward operation
func (k Keeper) RecordClaimReward(ctx sdk.Context, creator string, denom string, amount sdk.Int) error {
	assetProfile, found := k.apKeeper.GetEntry(ctx, denom)
	if !found {
		return sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return sdkerrors.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)
	if !found {
		return sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", creator)
	}

	// Subtract the withdrawn amount from the unclaimed balance
	err := commitments.SubRewardsUnclaimed(sdk.NewCoin(denom, amount))
	if err != nil {
		return err
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	err = k.HandleWithdrawFromCommitment(ctx, &commitments, addr, withdrawCoins, false)
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
