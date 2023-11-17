package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
)

// Vest converts user's commitment to vesting - start with unclaimed rewards and if it's not enough deduct from committed bucket
// mainly utilized for Eden
func (k msgServer) Vest(goCtx context.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.ProcessTokenVesting(ctx, msg.Denom, msg.Amount, msg.Creator); err != nil {
		return &types.MsgVestResponse{}, err
	}

	return &types.MsgVestResponse{}, nil
}

// Vesting token
// Check if vesting entity count is not exceeding the maximum and if it is fine, creates a new vesting entity
// Deduct from unclaimed bucket. If it is insufficent, deduct from committed bucket as well.
func (k Keeper) ProcessTokenVesting(ctx sdk.Context, denom string, amount sdk.Int, creator string) error {
	vestingInfo, _ := k.GetVestingInfo(ctx, denom)

	if vestingInfo == nil {
		return sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", denom)
	}

	commitments := k.GetCommitments(ctx, creator)

	// Create vesting tokens entry and add to commitments
	vestingTokens := commitments.GetVestingTokens()
	if vestingInfo.NumMaxVestings <= (int64)(len(vestingTokens)) {
		return sdkerrors.Wrapf(types.ErrExceedMaxVestings, "creator: %s", creator)
	}

	commitments, err := k.DeductClaimed(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	vestingTokens = append(vestingTokens, types.VestingTokens{
		Denom:           vestingInfo.VestingDenom,
		TotalAmount:     amount,
		UnvestedAmount:  amount,
		EpochIdentifier: vestingInfo.EpochIdentifier,
		NumEpochs:       vestingInfo.NumEpochs,
		CurrentEpoch:    0,
	})
	commitments.VestingTokens = vestingTokens

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, creator, sdk.Coins{sdk.NewCoin(denom, amount)})

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}
