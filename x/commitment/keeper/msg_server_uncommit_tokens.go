package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) UncommitTokens(goCtx context.Context, msg *types.MsgUncommitTokens) (*types.MsgUncommitTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	// Check if the committed tokens have enough amount to be uncommitted
	committedToken, _ := commitments.GetCommittedTokensForDenom(msg.Denom)
	// committedAmount := commitments.GetCommittedAmountForDenom(msg.Denom)

	if committedToken.Amount.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientCommittedTokens, "creator: %s, denom: %s", msg.Creator, msg.Denom)
	}

	// Update the committed tokens amount
	committedToken.Amount = committedToken.Amount.Sub(msg.Amount)

	// Update the uncommitted tokens amount
	uncommittedToken, found := commitments.GetUncommittedTokensForDenom(msg.Denom)

	if found {
		uncommittedToken.Amount = uncommittedToken.Amount.Add(msg.Amount)
	} else {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedTokens = append(uncommittedTokens, &types.UncommittedTokens{
			Denom:  msg.Denom,
			Amount: msg.Amount,
		})
		commitments.UncommittedTokens = uncommittedTokens
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.HookCommitmentChanged(ctx, msg.Creator, sdk.NewCoin(msg.Denom, msg.Amount))

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
