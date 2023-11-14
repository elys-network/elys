package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
)

// CancelVest cancel the user's vesting and the user reject to get vested tokens
func (k msgServer) CancelVest(goCtx context.Context, msg *types.MsgCancelVest) (*types.MsgCancelVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo, _ := k.GetVestingInfo(ctx, msg.Denom)
	if vestingInfo == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, msg.Creator)

	remainingToCancel := msg.Amount

	newVestingTokens := make([]*types.VestingTokens, 0, len(commitments.VestingTokens))

	for _, vesting := range commitments.VestingTokens {
		cancelAmount := sdk.MinInt(remainingToCancel, vesting.UnvestedAmount)
		vesting.TotalAmount = vesting.TotalAmount.Sub(cancelAmount)
		vesting.UnvestedAmount = vesting.UnvestedAmount.Sub(cancelAmount)

		if !vesting.TotalAmount.IsZero() {
			newVestingTokens = append(newVestingTokens, vesting)
		}

		remainingToCancel = remainingToCancel.Sub(cancelAmount)
		if remainingToCancel.IsZero() {
			break
		}
	}

	commitments.VestingTokens = newVestingTokens

	if !remainingToCancel.IsZero() {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientVestingTokens, "denom: %s, amount: %s", msg.Denom, msg.Amount)
	}

	// Update the unclaimed tokens amount
	commitments.AddRewardsUnclaimed(sdk.NewCoin(msg.Denom, msg.Amount))
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, msg.Creator, sdk.Coins{sdk.NewCoin(msg.Denom, msg.Amount)})

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgCancelVestResponse{}, nil
}
