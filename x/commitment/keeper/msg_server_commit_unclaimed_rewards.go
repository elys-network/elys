package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// CommitUnclaimedRewards commit the tokens on unclaimed store to committed
func (k msgServer) CommitUnclaimedRewards(goCtx context.Context, msg *types.MsgCommitUnclaimedRewards) (*types.MsgCommitUnclaimedRewardsResponse, error) {
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

	// Check if the unclaimed tokens have enough amount to be committed
	rewardUnclaimed, found := commitments.GetRewardsUnclaimedForDenom(msg.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientRewardsUnclaimed, "creator: %s", msg.Creator)
	}

	if rewardUnclaimed.Amount.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientRewardsUnclaimed, "creator: %s, denom: %s", msg.Creator, msg.Denom)
	}

	// Update the unclaimed tokens amount
	rewardUnclaimed.Amount = rewardUnclaimed.Amount.Sub(msg.Amount)

	// Update the committed tokens amount
	committedToken, found := commitments.GetCommittedTokensForDenom(msg.Denom)
	if found {
		committedToken.Amount = committedToken.Amount.Add(msg.Amount)
	} else {
		committedTokens := commitments.GetCommittedTokens()
		committedTokens = append(committedTokens, &types.CommittedTokens{
			Denom:  msg.Denom,
			Amount: msg.Amount,
		})
		commitments.CommittedTokens = committedTokens
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, msg.Creator, sdk.NewCoin(msg.Denom, msg.Amount))

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgCommitUnclaimedRewardsResponse{}, nil
}
