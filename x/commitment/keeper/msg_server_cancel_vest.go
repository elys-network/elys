package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

// CancelVest cancel the user's vesting and the user reject to get vested tokens
func (k msgServer) CancelVest(goCtx context.Context, msg *types.MsgCancelVest) (*types.MsgCancelVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.CancelVest(ctx, msg)
}

func (k Keeper) CancelVest(ctx sdk.Context, msg *types.MsgCancelVest) (*types.MsgCancelVestResponse, error) {
	if msg.Denom != ptypes.Eden {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	vestingInfo, _ := k.GetVestingInfo(ctx, ptypes.Eden)
	if vestingInfo == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", ptypes.Eden)
	}

	// claim pending rewards
	claimVestingMsg := types.MsgClaimVesting{Sender: msg.Creator}
	_, err := k.ClaimVesting(ctx, &claimVestingMsg)
	if err != nil {
		return nil, err
	}

	// Get the Commitments for the creator
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	commitments := k.GetCommitments(ctx, creator)

	remainingToCancel := msg.Amount
	totalCancelled := sdkmath.ZeroInt()

	for i := len(commitments.VestingTokens) - 1; i >= 0; i-- {
		vesting := commitments.VestingTokens[i]
		if vesting.Denom != ptypes.Elys || vesting.NumBlocks == 0 || vesting.TotalAmount.IsZero() {
			continue
		}

		// rewards claimed, so claimedAmount = vestedSoFar
		cancelAmount := sdkmath.MinInt(remainingToCancel, vesting.TotalAmount.Sub(vesting.ClaimedAmount))

		vesting.TotalAmount = vesting.TotalAmount.Sub(vesting.ClaimedAmount).Sub(cancelAmount)
		vesting.ClaimedAmount = sdkmath.ZeroInt()
		// remaining blocks for the new vesting amount
		vesting.NumBlocks = max(0, vesting.NumBlocks-(ctx.BlockHeight()-vesting.StartBlock))
		vesting.StartBlock = ctx.BlockHeight()
		vesting.VestStartedTimestamp = ctx.BlockTime().Unix()

		// Update the num epochs for the reduced amount
		commitments.VestingTokens[i] = vesting

		remainingToCancel = remainingToCancel.Sub(cancelAmount)
		totalCancelled = totalCancelled.Add(cancelAmount)
	}

	newVestingTokens := []*types.VestingTokens{}
	for _, vesting := range commitments.VestingTokens {
		if vesting.ClaimedAmount.GTE(vesting.TotalAmount) {
			continue
		}
		newVestingTokens = append(newVestingTokens, vesting)
	}

	commitments.VestingTokens = newVestingTokens

	if totalCancelled.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrInsufficientVestingTokens, "denom: %s, amount: %s", ptypes.Eden, msg.Amount)
	}
	ctx.Logger().Info("Successfully Cancelled vesting token",
		"creator", msg.Creator,
		"amount", totalCancelled.String(),
		"denom", msg.Denom)

	// Update the unclaimed tokens amount
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, totalCancelled))

	prev := k.GetTotalSupply(ctx)
	prev.TotalEdenSupply = prev.TotalEdenSupply.Add(totalCancelled)
	prev.TotalEdenVested = prev.TotalEdenVested.Sub(totalCancelled)
	k.SetTotalSupply(ctx, prev)
	k.SetCommitments(ctx, commitments)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeCancelledAmount, totalCancelled.String()),
			sdk.NewAttribute(types.AttributeDenom, ptypes.Eden),
		),
	)

	return &types.MsgCancelVestResponse{}, nil
}
