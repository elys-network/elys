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
	uncommittedTokens := commitments.GetUncommittedTokens()
	uncommittedAmount := sdk.NewInt(0)
	for _, token := range uncommittedTokens {
		if token.Denom == msg.Denom {
			uncommittedAmount = token.Amount
			break
		}
	}

	if uncommittedAmount.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientUncommittedTokens, "creator: %s, denom: %s", msg.Creator, msg.Denom)
	}

	// Update the uncommitted tokens amount
	for _, token := range uncommittedTokens {
		if token.Denom == msg.Denom {
			token.Amount = token.Amount.Sub(msg.Amount)
			break
		}
	}

	// Update the committed tokens amount
	committedTokens := commitments.GetCommittedTokens()
	found = false
	for _, token := range committedTokens {
		if token.Denom == msg.Denom {
			token.Amount = token.Amount.Add(msg.Amount)
			found = true
			break
		}
	}

	if !found {
		committedTokens = append(committedTokens, &types.CommittedTokens{
			Denom:  msg.Denom,
			Amount: msg.Amount,
		})
	}

	// Update the commitments
	commitments.CommittedTokens = committedTokens
	commitments.UncommittedTokens = uncommittedTokens
	k.SetCommitments(ctx, commitments)

	return &types.MsgCommitTokensResponse{}, nil
}
