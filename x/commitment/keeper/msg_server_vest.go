package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) Vest(goCtx context.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo := k.GetVestingInfo(ctx, msg.Denom)

	if vestingInfo == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	// Create vesting tokens entry and add to commitments
	vestingTokens := commitments.GetVestingTokens()
	if vestingInfo.NumMaxVestings <= (int64)(len(vestingTokens)) {
		return nil, sdkerrors.Wrapf(types.ErrExceedMaxVestings, "creator: %s", msg.Creator)
	}

	commitments, err := k.DeductCommitments(ctx, msg.Creator, msg.Denom, msg.Amount)
	if err != nil {
		return nil, err
	}

	vestingTokens = append(vestingTokens, &types.VestingTokens{
		Denom:           vestingInfo.VestingDenom,
		TotalAmount:     msg.Amount,
		UnvestedAmount:  msg.Amount,
		EpochIdentifier: vestingInfo.EpochIdentifier,
		NumEpochs:       vestingInfo.NumEpochs,
		CurrentEpoch:    0,
	})
	commitments.VestingTokens = vestingTokens

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

	return &types.MsgVestResponse{}, nil
}
