package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
)

func (k msgServer) Vest(goCtx context.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Denom != "eden" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	// Get user's uncommitted balance
	uncommittedToken, _ := commitments.GetUncommittedTokensForDenom(msg.Denom)

	requestedAmount := msg.Amount

	// Check if there are enough uncommitted tokens to withdraw
	if uncommittedToken.Amount.LT(requestedAmount) {
		// Calculate the difference between the requested amount and the available uncommitted balance
		difference := requestedAmount.Sub(uncommittedToken.Amount)

		committedToken, found := commitments.GetCommittedTokensForDenom(msg.Denom)
		if found {
			if committedToken.Amount.GTE(difference) {
				// Uncommit the required committed tokens
				committedToken.Amount = committedToken.Amount.Sub(difference)
				requestedAmount = requestedAmount.Sub(difference)
			} else {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "not enough tokens to withdraw")
			}
		}
	}

	// Subtract the withdrawn amount from the uncommitted balance
	uncommittedToken.Amount = uncommittedToken.Amount.Sub(requestedAmount)

	// Create vesting tokens entry and add to commitments
	vestingTokens := commitments.GetVestingTokens()
	vestingTokens = append(vestingTokens, &types.VestingTokens{
		Denom:           k.stakingKeeper.BondDenom(ctx), // TODO: param or map
		TotalAmount:     msg.Amount,
		UnvestedAmount:  msg.Amount,
		EpochIdentifier: epochstypes.DayEpochID, // TODO: gov param?
		NumEpochs:       180,
		CurrentEpoch:    0,
	})
	commitments.VestingTokens = vestingTokens

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	return &types.MsgVestResponse{}, nil
}
