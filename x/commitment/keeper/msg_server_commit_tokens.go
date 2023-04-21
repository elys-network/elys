package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) CommitTokens(goCtx context.Context, msg *types.MsgCommitTokens) (*types.MsgCommitTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	// Check if the uncommitted tokens have enough amount to be committed
	uncommittedToken, _ := commitments.GetUncommittedTokensForDenom(msg.Denom)

	if uncommittedToken.Amount.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientUncommittedTokens, "creator: %s, denom: %s", msg.Creator, msg.Denom)
	}

	// Update the uncommitted tokens amount
	uncommittedToken.Amount = uncommittedToken.Amount.Sub(msg.Amount)

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

	return &types.MsgCommitTokensResponse{}, nil
}
