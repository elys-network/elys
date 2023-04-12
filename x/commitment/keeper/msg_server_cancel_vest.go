package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) CancelVest(goCtx context.Context, msg *types.MsgCancelVest) (*types.MsgCancelVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo := k.GetVestingInfo(ctx, msg.Denom)

	if vestingInfo == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	remainingToCancel := msg.Amount

	for i := 0; i < len(commitments.VestingTokens) && !remainingToCancel.IsZero(); {
		vesting := commitments.VestingTokens[i]

		if vesting.Denom == msg.Denom {
			cancelAmount := sdk.MinInt(remainingToCancel, vesting.UnvestedAmount)

			vesting.TotalAmount = vesting.TotalAmount.Sub(cancelAmount)
			vesting.UnvestedAmount = vesting.UnvestedAmount.Sub(cancelAmount)

			if vesting.TotalAmount.IsZero() {
				commitments.VestingTokens = append(commitments.VestingTokens[:i], commitments.VestingTokens[i+1:]...)
			} else {
				i++
			}

			remainingToCancel = remainingToCancel.Sub(cancelAmount)
		} else {
			i++
		}
	}

	if !remainingToCancel.IsZero() {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientVestingTokens, "denom: %s, amount: %s", msg.Denom, msg.Amount)
	}

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

	k.SetCommitments(ctx, commitments)

	return &types.MsgCancelVestResponse{}, nil
}
