package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// UncommitTokens uncommits the tokens from committed store and make it liquid immediately.
func (k msgServer) UncommitTokens(goCtx context.Context, msg *types.MsgUncommitTokens) (*types.MsgUncommitTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	assetProfile, found := k.apKeeper.GetEntry(ctx, msg.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", msg.Denom)
	}

	if !assetProfile.CommitEnabled {
		return nil, sdkerrors.Wrapf(types.ErrCommitDisabled, "denom: %s", msg.Denom)
	}

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	// Deduct from committed tokens
	err := commitments.DeductFromCommitted(msg.Denom, msg.Amount, uint64(ctx.BlockTime().Unix()))
	if err != nil {
		return nil, err
	}
	k.SetCommitments(ctx, commitments)

	liquidCoins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))

	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// Send the coins to the user's account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, liquidCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to send liquid tokens")
	}

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, msg.Creator, sdk.Coins{sdk.NewCoin(msg.Denom, msg.Amount)})

	// Emit Hook if Eden is uncommitted
	if msg.Denom == ptypes.Eden {
		k.EdenUncommitted(ctx, msg.Creator, sdk.NewCoin(msg.Denom, msg.Amount))
	}

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgUncommitTokensResponse{}, nil
}
