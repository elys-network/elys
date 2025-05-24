package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

// Vest converts user's commitment to vesting - start with unclaimed rewards and if it's not enough deduct from committed bucket
// mainly utilized for Eden
func (k msgServer) Vest(goCtx context.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	err := k.ProcessTokenVesting(ctx, msg.Denom, msg.Amount, creator)
	if err != nil {
		return &types.MsgVestResponse{}, err
	}

	return &types.MsgVestResponse{}, nil
}

// Vesting token
// Check if vesting entity count is not exceeding the maximum and if it is fine, creates a new vesting entity
// Deduct from unclaimed bucket. If it is insufficient, deduct from committed bucket as well.
func (k Keeper) ProcessTokenVesting(ctx sdk.Context, denom string, amount math.Int, creator sdk.AccAddress) error {
	vestingInfo, _ := k.GetVestingInfo(ctx, denom)

	if vestingInfo == nil {
		return errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", denom)
	}

	commitments := k.GetCommitments(ctx, creator)

	// Create vesting tokens entry and add to commitments
	vestingTokens := commitments.GetVestingTokens()
	if vestingInfo.NumMaxVestings <= (int64)(len(vestingTokens)) {
		return errorsmod.Wrapf(types.ErrExceedMaxVestings, "creator: %s", creator.String())
	}

	commitments, err := k.DeductClaimed(ctx, creator, denom, amount)
	if err != nil {
		return err
	}

	vestingTokens = append(vestingTokens, &types.VestingTokens{
		Denom:                vestingInfo.VestingDenom,
		TotalAmount:          amount,
		ClaimedAmount:        math.ZeroInt(),
		StartBlock:           ctx.BlockHeight(),
		NumBlocks:            vestingInfo.NumBlocks,
		VestStartedTimestamp: ctx.BlockTime().Unix(),
	})
	commitments.VestingTokens = vestingTokens

	if denom == ptypes.Eden {
		prev := k.GetTotalSupply(ctx)
		prev.TotalEdenSupply = prev.TotalEdenSupply.Sub(amount)
		prev.TotalEdenVested = prev.TotalEdenVested.Add(amount)
		k.SetTotalSupply(ctx, prev)
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, creator.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)

	return nil
}
