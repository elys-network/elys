package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// CommitClaimedRewards commit the tokens on unclaimed store to committed
func (k msgServer) CommitClaimedRewards(goCtx context.Context, msg *types.MsgCommitClaimedRewards) (*types.MsgCommitClaimedRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, msg.Denom)
	if !found {
		return nil, errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", msg.Denom)
	}

	if !assetProfile.CommitEnabled {
		return nil, errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", msg.Denom)
	}

	params := k.GetParams(ctx)
	params.TotalCommitted = params.TotalCommitted.Add(sdk.NewCoin(msg.Denom, msg.Amount))
	k.SetParams(ctx, params)

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, msg.Creator)

	// Decrease unclaimed tokens amount
	err := commitments.SubClaimed(sdk.NewCoin(msg.Denom, msg.Amount))
	if err != nil {
		return nil, err
	}

	// Increase committed tokens
	commitments.AddCommittedTokens(msg.Denom, msg.Amount, uint64(ctx.BlockTime().Unix()))

	// Update the commitments
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

	return &types.MsgCommitClaimedRewardsResponse{}, nil
}
