package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

// CommitClaimedRewards commit the tokens on unclaimed store to committed
func (k msgServer) CommitClaimedRewards(goCtx context.Context, msg *types.MsgCommitClaimedRewards) (*types.MsgCommitClaimedRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Creator)
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, msg.Denom)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "denom: %s", msg.Denom)
	}

	if !assetProfile.CommitEnabled {
		return nil, errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", msg.Denom)
	}

	params := k.GetParams(ctx)
	params.TotalCommitted = params.TotalCommitted.Add(sdk.NewCoin(msg.Denom, msg.Amount))
	k.SetParams(ctx, params)

	// Get the Commitments for the creator
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	commitments := k.GetCommitments(ctx, creator)

	if msg.Denom == ptypes.Eden {
		if commitments.GetCommittedAmountForDenom(ptypes.Eden).IsPositive() {
			err := k.hooks.BeforeEdenCommitChange(ctx, sender)
			if err != nil {
				return nil, err
			}
		} else {
			err := k.hooks.BeforeEdenInitialCommit(ctx, sender)
			if err != nil {
				return nil, err
			}
		}
	}

	if msg.Denom == ptypes.EdenB {
		if commitments.GetCommittedAmountForDenom(ptypes.EdenB).IsPositive() {
			err := k.hooks.BeforeEdenBCommitChange(ctx, sender)
			if err != nil {
				return nil, err
			}
		} else {
			err := k.hooks.BeforeEdenBInitialCommit(ctx, sender)
			if err != nil {
				return nil, err
			}
		}
	}

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
	err = k.CommitmentChanged(ctx, creator, sdk.Coins{sdk.NewCoin(msg.Denom, msg.Amount)})
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

	return &types.MsgCommitClaimedRewardsResponse{}, nil
}
