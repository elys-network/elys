package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// Update commitments for claim reward operation
func (k Keeper) RecordClaimReward(ctx sdk.Context, creator string, denom string, amount math.Int, withdrawMode types.EarnType) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return errorsmod.Wrapf(types.ErrWithdrawDisabled, "denom: %s", denom)
	}

	// uses asset profile denom
	denom = assetProfile.Denom

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, creator)
	if !found {
		return errorsmod.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", creator)
	}

	// Withdraw reward not depending on program type
	switch withdrawMode {
	case types.EarnType_ALL_PROGRAM:
		// Subtract the withdrawn amount from the unclaimed balance
		err := commitments.SubRewardsUnclaimed(sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	case types.EarnType_USDC_PROGRAM:
		// Subtract the withdrawn amount from the unclaimed balance
		err := commitments.SubRewardsUnclaimedForUSDCDeposit(sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	case types.EarnType_ELYS_PROGRAM:
		// Subtract the withdrawn amount from the unclaimed balance
		err := commitments.SubRewardsUnclaimedForElysStaking(sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	case types.EarnType_EDEN_PROGRAM:
		// Subtract the withdrawn amount from the unclaimed balance
		err := commitments.SubRewardsUnclaimedForEdenCommitted(sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	case types.EarnType_EDENB_PROGRAM:
		// Subtract the withdrawn amount from the unclaimed balance
		err := commitments.SubRewardsUnclaimedForEdenBCommitted(sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	default:
		return errorsmod.Wrapf(types.ErrUnsupportedWithdrawMode, "creator: %s", creator)
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
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
