package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
)

// CancelVest cancel the user's vesting and the user reject to get vested tokens
func (k msgServer) CancelVest(goCtx context.Context, msg *types.MsgCancelVest) (*types.MsgCancelVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo, _ := k.GetVestingInfo(ctx, msg.Denom)
	if vestingInfo == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, msg.Creator)

	remainingToCancel := msg.Amount

	remainVestingTokens := make([]*types.VestingTokens, 0, len(commitments.VestingTokens))
	newVestingTokens := make([]*types.VestingTokens, 0, len(commitments.VestingTokens))

	for i := len(commitments.VestingTokens) - 1; i >= 0; i-- {
		vesting := commitments.VestingTokens[i]
		cancelAmount := sdk.MinInt(remainingToCancel, vesting.UnvestedAmount)
		if vesting.NumEpochs == 0 || vesting.TotalAmount.IsZero() {
			continue
		}
		amtPerEpoch := vesting.TotalAmount.Quo(sdk.NewInt(vesting.NumEpochs))
		vesting.TotalAmount = vesting.TotalAmount.Sub(cancelAmount)
		vesting.UnvestedAmount = vesting.UnvestedAmount.Sub(cancelAmount)

		if amtPerEpoch.IsZero() {
			continue
		}
		// Update the num epochs for the reduced amount
		vesting.NumEpochs = vesting.TotalAmount.Quo(amtPerEpoch).Int64()

		if !vesting.TotalAmount.IsZero() {
			remainVestingTokens = append(remainVestingTokens, vesting)
		}

		remainingToCancel = remainingToCancel.Sub(cancelAmount)
		if remainingToCancel.IsZero() {
			newVestingTokens = append(newVestingTokens, commitments.VestingTokens[:i]...)

			break
		}
	}

	// Add the remaining vesting in reverse order.
	for i := len(remainVestingTokens) - 1; i >= 0; i-- {
		newVestingTokens = append(newVestingTokens, remainVestingTokens[i])
	}

	commitments.VestingTokens = newVestingTokens

	if !remainingToCancel.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrInsufficientVestingTokens, "denom: %s, amount: %s", msg.Denom, msg.Amount)
	}

	// Update the unclaimed tokens amount
	commitments.AddClaimed(sdk.NewCoin(msg.Denom, msg.Amount))
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	err := k.AfterCommitmentChange(ctx, msg.Creator, sdk.Coins{sdk.NewCoin(msg.Denom, msg.Amount)})
	if err != nil {
		return nil, err
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

	return &types.MsgCancelVestResponse{}, nil
}
