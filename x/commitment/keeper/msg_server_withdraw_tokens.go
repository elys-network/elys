package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) WithdrawTokens(goCtx context.Context, msg *types.MsgWithdrawTokens) (*types.MsgWithdrawTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
	// Update the commitments
	k.SetCommitments(ctx, commitments)

	withdrawCoins := sdk.NewCoins(sdk.NewCoin(msg.Denom, requestedAmount))

	// Mint the withdrawn tokens to the user's account
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, withdrawCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to mint withdrawn tokens")
	}
	// Send the minted coins to the user's account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(msg.Creator), withdrawCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to send withdrawn tokens")
	}

	return &types.MsgWithdrawTokensResponse{}, nil
}
