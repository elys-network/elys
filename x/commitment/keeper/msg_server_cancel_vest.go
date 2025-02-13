package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// CancelVest cancel the user's vesting and the user reject to get vested tokens
func (k msgServer) CancelVest(goCtx context.Context, msg *types.MsgCancelVest) (*types.MsgCancelVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Denom != ptypes.Eden {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	vestingInfo, _ := k.GetVestingInfo(ctx, ptypes.Eden)
	if vestingInfo == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", ptypes.Eden)
	}

	// Get the Commitments for the creator
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	commitments := k.GetCommitments(ctx, creator)

	remainingToCancel := msg.Amount

	for i := len(commitments.VestingTokens) - 1; i >= 0; i-- {
		vesting := commitments.VestingTokens[i]
		if vesting.Denom != ptypes.Elys || vesting.NumBlocks == 0 || vesting.TotalAmount.IsZero() {
			continue
		}
		cancelAmount := sdkmath.MinInt(remainingToCancel, vesting.TotalAmount.Sub(vesting.ClaimedAmount))
		vesting.TotalAmount = vesting.TotalAmount.Sub(cancelAmount)
		// Update the num epochs for the reduced amount
		commitments.VestingTokens[i] = vesting

		remainingToCancel = remainingToCancel.Sub(cancelAmount)
	}

	newVestingTokens := []*types.VestingTokens{}
	for _, vesting := range commitments.VestingTokens {
		if vesting.ClaimedAmount.GTE(vesting.TotalAmount) {
			continue
		}
		newVestingTokens = append(newVestingTokens, vesting)
	}

	commitments.VestingTokens = newVestingTokens

	if !remainingToCancel.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrInsufficientVestingTokens, "denom: %s, amount: %s", ptypes.Eden, msg.Amount)
	}

	// Update the unclaimed tokens amount
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, msg.Amount))

	prev := k.GetTotalSupply(ctx)
	prev.TotalEdenSupply = prev.TotalEdenSupply.Add(msg.Amount)
	k.SetTotalSupply(ctx, prev)
	k.SetCommitments(ctx, commitments)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, ptypes.Eden),
		),
	)

	return &types.MsgCancelVestResponse{}, nil
}
